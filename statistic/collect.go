package statistic

import (
	"context"
	"log"

	"github.com/andriiyaremenko/mg/dto"
)

func StartCollection(
	ctx context.Context,
	log *log.Logger,
	statEndpoint string,
) *Agent {
	log.Println("starting statistic collection agent")

	agent := &Agent{
		sc:      &Collector{stats: make(map[int]dto.TargetStatistic)},
		log:     log,
		statsCh: make(chan []dto.TargetStatistic, 1),
		done:    make(chan struct{}),
	}

	agent.start(ctx, statEndpoint)

	return agent
}
