package handler

import (
	"context"

	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
)

func HandleTargetError(numWorkers int) pipelines.Handler[dto.TargetResponse, dto.Statistic] {
	return &pipelines.BaseHandler[dto.TargetResponse, dto.Statistic]{
		NWorkers: numWorkers,
		HandleFunc: func(ctx context.Context, w pipelines.EventWriter[dto.Statistic], e pipelines.Event[dto.TargetResponse]) {
			target := e.Payload
			if target.ID != 0 {

				w.Write(pipelines.Event[dto.Statistic]{
					Payload: dto.Statistic{
						Success: 0,
						Error:   1,
					},
				})

				return
			}

			w.Write(pipelines.NewErr[dto.Statistic](e.Err))
		},
	}
}
