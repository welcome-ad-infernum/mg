package handler

import (
	"context"

	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
)

func HandleTargetError(ctx context.Context, w pipelines.EventWriter[dto.Statistic], err error) {
	if _, ok := err.(*pipelines.Error[bool]); ok {
		w.Write(pipelines.Event[dto.Statistic]{
			Payload: dto.Statistic{
				Success: 0,
				Error:   1,
			},
		})

		return
	}

	w.Write(pipelines.NewErrEvent[dto.Statistic](err))
}
