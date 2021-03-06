package handler

import (
	"context"
	"log"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/dto"
	"github.com/andriiyaremenko/pipelines"
)

func ProcessTargetResponse(log *log.Logger) pipelines.Handler[dto.TargetResponse, dto.Statistic] {
	return func(ctx context.Context, w pipelines.EventWriter[dto.Statistic], target dto.TargetResponse) {
		if target.Code < 500 {
			log.Println(
				logw.Debug.WithString("target", target.URL),
				"target is alive",
			)
			w.Write(pipelines.Event[dto.Statistic]{
				Payload: dto.Statistic{
					Success: 1,
					Error:   0,
				},
			})

			return
		}

		log.Println(
			logw.Debug.
				WithString("target", target.URL).
				WithInt("status_code", target.Code),
			"target returned error",
		)

		w.Write(pipelines.Event[dto.Statistic]{
			Payload: dto.Statistic{
				Success: 0,
				Error:   1,
			},
		})
	}
}
