package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	d "ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"time"
)

var logger = log.New(os.Stderr, "[Parsing] ", log.Lmsgprefix|log.Ltime)

type Collector struct {
	processUrls []string
}

func (c *Collector) CollectUrls(startDate time.Time) {
	currentDate := startDate
	endDate := time.Now().AddDate(0, 0, 1)

	for currentDate.Format("2006-01-02") != endDate.Format("2006-01-02") {
		c.trySaveUrl(&currentDate)
		currentDate = currentDate.AddDate(0, 0, 1)
	}
}

func (c *Collector) trySaveUrl(date *time.Time) {
	urls := r.FindAllByDate(date)
	if len(urls) > 0 {
		c.processUrls = append(c.processUrls, urls...)
		return
	}

	dateUrl := dateUrl(date)
	logger.Println(dateUrl)
	roundId := roundIds(dateUrl)

	for _, v := range roundId {
		url := statUrl(date, v)
		c.processUrls = append(c.processUrls, url)
		r.SaveDate(&d.Link{Link: url, Date: date.Format("2006-01-02")})
	}
}

func roundIds(url string) []string {
	resp, err := http.Get(url)

	if err != nil {
		logger.Println(err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		logger.Println(resp.Status)
		return nil
	}

	var rounds []d.RoundInDate
	json.NewDecoder(resp.Body).Decode(&rounds)

	roundsIds := make([]string, 0, cap(rounds))
	for _, v := range rounds {
		roundsIds = append(roundsIds, v.Round)
	}
	return roundsIds
}
