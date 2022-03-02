package handler

import (
	"context"
	"log"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/tinycqs/command"
)

func HandleErrors(log *log.Logger) command.Handler {
	return &command.BaseHandler{
		Type: command.CatchAllErrorEventType,
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			errEvent, ok := e.(*command.ErrEvent)
			if !ok {
				log.Println(logw.Error, e.Err())
				return
			}

			log.Println(logw.Error, errEvent.Unwrap())
		},
	}
}
