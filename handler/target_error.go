package handler

import (
	"context"
	"encoding/json"
	"log"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func TargetError(log *log.Logger) command.Handler {
	return &command.BaseHandler{
		Type: "TARGET_ERROR",
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			target := new(dto.TargetError)
			if err := json.Unmarshal(e.Payload(), target); err != nil {
				w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target record")))
				return
			}

			log.Println(
				logw.Warn.
					WithString("target", target.URL).
					WithInt("status_code", target.ErrCode),
				"target returned error",
			)
		},
	}
}
