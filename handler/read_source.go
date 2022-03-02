package handler

import (
	"context"
	"encoding/json"

	"github.com/andriiyaremenko/mg/source"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func ReadSource(s source.Source) command.Handler {
	return &command.BaseHandler{
		Type: "READ_SOURCE",
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			for target, keep, err := s(); keep; {
				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "error:")))

					continue
				}

				b, err := json.Marshal(target)
				if err != nil {
					w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target record")))

					continue
				}

				w.Write(command.E{
					Type: "NUKE_TARGET",
					P:    b,
				})
			}
		},
	}
}
