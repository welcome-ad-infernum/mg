package handler

import (
	"context"
	"encoding/json"

	"github.com/andriiyaremenko/mg/source"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func ReadSource(s source.Source, amountWorkers int) command.Handler {
	return &command.BaseHandler{
		Type:     "READ_SOURCE",
		NWorkers: 1,
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			for target, keep, err := s(); keep; target, keep, err = s() {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "error:")))

					continue
				}

				b, err := json.Marshal(target)
				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target record")))

					continue
				}

				for i := amountWorkers; i > 0; i-- {
					w.Write(command.E{
						Type: "NUKE_TARGET",
						P:    b,
					})
				}
			}
		},
	}
}
