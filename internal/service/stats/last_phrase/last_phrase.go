package last_phrase

import (
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/utils"
)

const (
	lastPhrasesLimit = 500
)

type lastPhrase struct {
	Name        string
	RoundID     int32
	Phrase      string
	TimeOfDeath string
}

var (
	lastPhrasesPool = make([]lastPhrase, 0, lastPhrasesLimit)
)

func fillPoll() {
	r.Database.Raw(`SELECT *
						FROM (
							SELECT distinct on (last_phrase, real_name) last_phrase as phrase, real_name as name, root_id as round_id, time_of_death as Time_of_death
							from deaths
							WHERE last_phrase <> '') as d
						ORDER BY random()
						LIMIT ?
		`, lastPhrasesLimit).
		Scan(&lastPhrasesPool)
}

func GetRandomLastPhrase() lastPhrase {
	if len(lastPhrasesPool) == 0 {
		fillPoll()
	}
	lastIndex := len(lastPhrasesPool) - 1
	phrase := lastPhrasesPool[lastIndex]
	lastPhrasesPool = utils.RemoveByIndex(lastPhrasesPool, lastIndex)
	return phrase
}
