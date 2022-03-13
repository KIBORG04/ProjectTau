package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	d "ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	u "ssstatistics/internal/utils"
	"time"
)

var loggerStats = log.New(os.Stderr, "[Statistics] ", log.Lmsgprefix|log.Ltime)
var loggerGet = log.New(os.Stderr, "[GET Request] ", log.Lmsgprefix|log.Ltime)

type Collector struct {
	processUrls []string
	Logs        []string
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
	if urls == nil {
		return
	}
	if len(urls) > 0 {
		c.processUrls = append(c.processUrls, urls...)
		return
	}

	dateUrl := dateUrl(date)
	roundId := c.roundIds(dateUrl)
	if roundId == nil {
		return
	}

	for _, v := range roundId {
		url := statUrl(date, v)
		c.processUrls = append(c.processUrls, url)
		r.Save(&d.Link{Link: url, Date: date.Format("2006-01-02")})
	}
}

func (c *Collector) CollectStatistics() {
	for _, url := range c.processUrls {
		c.collectByUrl(url)
	}
}

func (c *Collector) saveLogs(logger *log.Logger, text interface{}) {
	text_str := fmt.Sprintln(text)
	c.Logs = append(c.Logs, text_str)
	logger.Println(text)
}

func (c *Collector) requestGET(url string) *http.Response {
	resp, err := http.Get(url)

	c.saveLogs(loggerGet, url)
	if err != nil {
		c.saveLogs(loggerGet, err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		c.saveLogs(loggerGet, resp.Status)
		return nil
	}

	return resp
}

func (c *Collector) collectByUrl(url string) {
	roundId := u.RoundId.FindString(url)
	if len(roundId) == 0 {
		c.saveLogs(loggerStats, fmt.Sprintf("%s not contain digits of the round", url))
		return
	}

	_, err := r.FindByRoundId(roundId)
	if err == nil {
		return
	}

	resp := c.requestGET(url)

	var root d.Root
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	dec.Decode(&root)

	r.Save(&root)
}

func (c *Collector) roundIds(url string) []string {
	resp := c.requestGET(url)
	if resp == nil {
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
