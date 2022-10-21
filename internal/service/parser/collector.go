package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	d "ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
	u "ssstatistics/internal/utils"
	"sync"
	"time"
)

type RoundDto struct {
	link string
	date string
}

var loggerStats = log.New(os.Stderr, "[Statistics] ", log.Lmsgprefix|log.Ltime)
var loggerGet = log.New(os.Stderr, "[GET Request] ", log.Lmsgprefix|log.Ltime)

type Collector struct {
	processUrls []*RoundDto
}

var (
	waitGroup sync.WaitGroup
	mutex     sync.Mutex
)

func RunRoundCollector() []string {
	startDate, _ := time.Parse("2006-01-02", stats.CurrentStatisticsDate)

	collector := Collector{}
	collector.CollectUrls(startDate)
	collector.CollectStatistics()
	return []string{
		fmt.Sprintf("%d Rounds Collected", len(collector.processUrls)),
	}
}

func (c *Collector) CollectUrls(startDate time.Time) {
	currentDate := startDate
	endDate := time.Now().AddDate(0, 0, 1)

	for currentDate.Format("2006-01-02") != endDate.Format("2006-01-02") {
		waitGroup.Add(1)

		go func(currentDate time.Time) {
			defer waitGroup.Done()
			c.trySaveUrl(&currentDate)
		}(currentDate)

		currentDate = currentDate.AddDate(0, 0, 1)
	}
	waitGroup.Wait()
}

func (c *Collector) trySaveUrl(date *time.Time) {
	dateUrl := dateUrl(date)
	roundIds := c.roundIds(dateUrl)
	if roundIds == nil {
		return
	}

	mutex.Lock()
	for _, v := range roundIds {
		roundId := u.RoundId.FindString(v)
		if len(roundId) == 0 {
			c.saveLogs(loggerStats, fmt.Sprintf("%s not contain digits of the round", v))
			return
		}
		_, err := r.FindByRoundId(roundId)
		if err == nil {
			continue
		}
		url := statUrl(date, v)
		c.processUrls = append(c.processUrls, &RoundDto{
			link: url,
			date: date.Format("2006-01-02"),
		})
	}
	mutex.Unlock()
}

func (c *Collector) CollectStatistics() {
	for _, link := range c.processUrls {
		waitGroup.Add(1)
		go func(link *RoundDto) {
			defer waitGroup.Done()
			c.collectByUrl(link)
		}(link)
	}
	waitGroup.Wait()
}

func (c *Collector) saveLogs(logger *log.Logger, text any) {
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

func (c *Collector) collectByUrl(link *RoundDto) {
	resp := c.requestGET(link.link)

	var root d.Root
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&root)
	if err != nil {
		fmt.Println(err)
	}
	root.Date = link.date

	mutex.Lock()
	r.Save(&root)
	mutex.Unlock()
}

func (c *Collector) roundIds(url string) []string {
	resp := c.requestGET(url)
	if resp == nil {
		return nil
	}

	var rounds []d.RoundInDate
	err := json.NewDecoder(resp.Body).Decode(&rounds)
	if err != nil {
		return nil
	}

	roundsIds := make([]string, 0, cap(rounds))
	for _, v := range rounds {
		roundsIds = append(roundsIds, v.Round)
	}

	return roundsIds
}
