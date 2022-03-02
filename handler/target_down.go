package handler

import (
	"context"
	"encoding/json"
	"log"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/logwriter/color"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/tinycqs/command"
	"github.com/pkg/errors"
)

func TargetDown(log *log.Logger) command.Handler {
	return &command.BaseHandler{
		Type: "TARGET_DOWN",
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			target := new(dto.Target)
			if err := json.Unmarshal(e.Payload(), target); err != nil {
				w.Write(command.NewErrEvent(e, errors.Wrap(err, "bad target record")))
				return
			}

			log.Println(
				logw.Info.WithString("target", target.URL),
				color.ColorizeText(color.ANSIColorRed, "target is down"),
			)
		},
	}
}
