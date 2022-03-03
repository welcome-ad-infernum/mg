package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func CollectStatistic(ctx context.Context, log *log.Logger, statEndpoint string) command.Handler {
	ch := make(chan dto.TargetStatistic)

	go startStatisticCollection(ctx, log, statEndpoint, ch)

	return &command.BaseHandler{
		Type:     "COLLECT_STATISTIC",
		NWorkers: 1,
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			var stats dto.TargetStatistic
			if err := json.Unmarshal(e.Payload(), &stats); err != nil {
				w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target statistic record")))

				return
			}

			ch <- stats
		},
	}
}

func startStatisticCollection(ctx context.Context, log *log.Logger, statEndpoint string, ch chan dto.TargetStatistic) {
	log.Println("starting statistic collection agent")

	cl := client.New()
	sc := &statisticCollector{stats: make(map[int]dto.Statistic)}
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Println("shutting down statistic collection agent")

			return
		case stat := <-ch:
			sc.Append(stat)
		case <-ticker.C:
			for target, stat := range sc.Get() {
				statEndpoint := fmt.Sprintf(statEndpoint, target)

				log.Println(logw.Info.WithMessage("sending statistic to %s", statEndpoint))

				b, err := json.Marshal(stat)
				if err != nil {
					log.Println(logw.Error, err)
					continue
				}

				r, err := cl.Post(statEndpoint, "application/json", bytes.NewReader(b))
				if err != nil {
					log.Println(logw.Error, err)
					continue
				}

				body := "<no request body>"
				b, err = io.ReadAll(r.Body)

				if err == nil {
					body = string(b)
				}

				r.Body.Close()

				if r.StatusCode >= 300 {
					log.Println(logw.Error, errors.Errorf("failed to send statistics: code %s, body: %s", r.Status, body))
					continue
				}

				log.Printf("statistic was sent, response status code %s", r.Status)
			}
		}
	}
}

type statisticCollector struct {
	rwMu sync.RWMutex

	stats map[int]dto.Statistic
}

func (sc *statisticCollector) Get() map[int]dto.Statistic {
	sc.rwMu.RLock()
	defer sc.rwMu.RUnlock()

	result := sc.stats
	sc.stats = make(map[int]dto.Statistic)

	return result
}

func (sc *statisticCollector) Append(stat dto.TargetStatistic) {
	sc.rwMu.Lock()
	defer sc.rwMu.Unlock()

	old, ok := sc.stats[stat.TargetID]
	if !ok {
		sc.stats[stat.TargetID] = stat.Statistic
		return
	}

	newStat := dto.Statistic{
		AgentUID: stat.AgentUID,
		Success:  old.Success + stat.Success,
		Error:    old.Error + stat.Error,
	}

	sc.stats[stat.TargetID] = newStat
}
