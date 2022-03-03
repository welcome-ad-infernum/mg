package handler

import (
	"context"

	"github.com/andriiyaremenko/tinycqs/command"
)

func LaunchAttack(amountWorkers int) command.Handler {
	return &command.BaseHandler{
		Type:     "LAUNCH_ATTACK",
		NWorkers: 1,
		HandleFunc: func(ctx context.Context, w command.EventWriter, e command.Event) {
			defer w.Done()

			for i := amountWorkers; i > 0; i-- {
				w.Write(command.E{
					Type: "NUKE_TARGET",
					P:    e.Payload(),
				})
			}
		},
	}
}
