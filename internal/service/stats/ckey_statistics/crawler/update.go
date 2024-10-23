package crawler

import (
	r "ssstatistics/internal/repository"
)

func SecretlyUpdateSomePlayers() []string {
	var ckeyUpdatedTimes []*struct {
		Ckey  string
		Count uint
	}
	r.Database.Raw(`
	select cs.player_id as Ckey, count(distinct (cs.crawler_updated_at)) as count
	from crawler_stats cs
	group by cs.player_id
	having count(distinct (cs.crawler_updated_at)) > 7
	limit 3;
	`).Scan(&ckeyUpdatedTimes)

	output := "Player crawler updated for: "
	for _, player := range ckeyUpdatedTimes {
		output += player.Ckey
		FetchPlayerStats(player.Ckey)
	}
	return []string{output}
}
