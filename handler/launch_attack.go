package handler

import (
	"context"

	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
)

func LaunchAttack(amountWorkers int) pipelines.Handler[dto.Target, dto.Target] {
	return &pipelines.BaseHandler[dto.Target, dto.Target]{
		NWorkers: 1,
		HandleFunc: func(ctx context.Context, w pipelines.EventWriter[dto.Target], e pipelines.Event[dto.Target]) {
			for i := amountWorkers; i > 0; i-- {
				w.Write(e)
			}
		},
	}
}
