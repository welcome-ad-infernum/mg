package statistic

import (
	"sync"

	"github.com/andriiyaremenko/mg/dto"
)

type Collector struct {
	rwMu sync.Mutex

	stats map[int]dto.TargetStatistic
}

func (sc *Collector) Unload() map[int]dto.TargetStatistic {
	sc.rwMu.Lock()
	defer sc.rwMu.Unlock()

	result := sc.stats
	sc.stats = make(map[int]dto.TargetStatistic)

	return result
}

func (sc *Collector) Append(stats ...dto.TargetStatistic) {
	sc.rwMu.Lock()
	defer sc.rwMu.Unlock()

	for _, stat := range stats {
		old, ok := sc.stats[stat.TargetID]
		if !ok {
			sc.stats[stat.TargetID] = stat

			continue
		}

		newStat := dto.TargetStatistic{
			AgentUID: stat.AgentUID,
			TargetID: stat.TargetID,
			Statistic: dto.Statistic{
				Success: old.Success + stat.Success,
				Error:   old.Error + stat.Error,
			},
		}

		sc.stats[stat.TargetID] = newStat
	}
}
