package statistic

import (
	"context"
	"log"
	"net/http"

	"github.com/andriiyaremenko/mg/dto"
)

func StartCollection(
	ctx context.Context,
	cl *http.Client,
	log *log.Logger,
	statEndpoint string,
) *Agent {
	log.Println("starting statistic collection agent")

	agent := &Agent{
		sc:      &Collector{stats: make(map[int]dto.TargetStatistic)},
		cl:      cl,
		log:     log,
		statsCh: make(chan []dto.TargetStatistic, 1),
		done:    make(chan struct{}),
	}

	agent.start(ctx, statEndpoint)

	return agent
}
