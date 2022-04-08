package handler

import (
	"context"

	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
)

func LaunchAttack(amountWorkers int) pipelines.Handler[dto.Target, dto.Target] {
	return func(ctx context.Context, w pipelines.EventWriter[dto.Target], target dto.Target) {
		for i := amountWorkers; i > 0; i-- {
			w.Write(pipelines.Event[dto.Target]{Payload: target})
		}
	}
}
