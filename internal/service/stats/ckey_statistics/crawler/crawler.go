package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"ssstatistics/internal/config"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/repository"
	"strconv"
	"strings"
	"time"
)

const (
	Url = "https://crawler.ss13.su/api/?ckey="
	// 24 hours is 1 day
	expirationStatsHours = 0
)

type crawlerResponse struct {
	Servername string `json:"servername,omitempty"`
	Count      string `json:"count,omitempty"`
}

var loggerGet = log.New(os.Stderr, "[CrawlerParser] ", log.Lmsgprefix|log.Ltime)

func FetchPlayerStats(ckey string) *domain.Player {
	if slices.Contains(config.Config.Secret.CrawlBlacklist, ckey) {
		return nil
	}

	player := repository.GetPlayer(ckey)
	if player == nil {
		loggerGet.Println("player not exist in database")
		return nil
	}
	if !isActualCrawlerStats(player) {
		loggerGet.Println("player has no crawler stats")
		err := setCrawlerStats(player)
		if err != nil {
			loggerGet.Println(err)
			return nil
		}
		repository.SavePlayer(player)
	}
	return player
}

func isActualCrawlerStats(player *domain.Player) bool {
	if len(player.CrawlerStats) == 0 {
		return false
	}
	now := time.Now()
	if now.Sub(player.CrawlerUpdatedAt).Hours() > expirationStatsHours {
		return false
	}
	return true
}

func setCrawlerStats(player *domain.Player) error {
	response, err := requestGET(player.Ckey)
	if err != nil {
		return err
	}
	var crawlerStats []domain.CrawlerStat
	for _, resp := range response {
		mins, _ := strconv.Atoi(resp.Count)
		crawlerStat := domain.CrawlerStat{
			PlayerID:         player.Ckey,
			ServerName:       resp.Servername,
			CrawlerUpdatedAt: time.Now().Truncate(24 * time.Hour),
			Minutes:          uint(mins),
		}
		crawlerStats = append(crawlerStats, crawlerStat)
	}

	crawlerStats = cutServerNames(crawlerStats)

	diff := getChangedResponses(player.CrawlerStats, crawlerStats)
	player.CrawlerStats = append(player.CrawlerStats, diff...)
	return nil
}

func requestGET(ckey string) ([]crawlerResponse, error) {
	resp, err := http.Get(Url + ckey)

	loggerGet.Println(ckey)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var stats []crawlerResponse
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		return nil, fmt.Errorf("json: %s\nparse error: %s", resp.Body, err)
	}

	// skip first useless data
	if stats[0].Servername == "" && stats[0].Count == "" {
		stats = stats[1:]
	}

	return stats, nil
}

func cutServerNames(servers []domain.CrawlerStat) []domain.CrawlerStat {
	var newCrawlerStat []domain.CrawlerStat
	sumCleanServersMinutes := make(map[string]*domain.CrawlerStat)

	for _, server := range servers {
		substr := containsSubstring(server.ServerName, config.Config.Secret.CorrectServerNames)
		if substr == "" {
			newCrawlerStat = append(newCrawlerStat, server)
			continue
		}
		serverToUpdate, exists := sumCleanServersMinutes[substr]
		if !exists {
			serverToUpdate = &server
			serverToUpdate.ServerName = substr
			sumCleanServersMinutes[substr] = serverToUpdate
		} else {
			serverToUpdate.Minutes += server.Minutes
		}
	}

	for _, server := range sumCleanServersMinutes {
		newCrawlerStat = append(newCrawlerStat, *server)
	}

	return newCrawlerStat
}

// returns substring
func containsSubstring(server string, correctServers []string) string {
	for _, serverName := range correctServers {
		if strings.Contains(server, serverName) {
			return serverName
		}
	}
	return ""
}

// return new array with new elements and changed minutes
func getChangedResponses(oldStats, newStats []domain.CrawlerStat) []domain.CrawlerStat {
	oldMap := make(map[string][]interface{})
	for _, oldStat := range oldStats {
		info := make([]interface{}, 2)
		info[0] = oldStat.Minutes
		info[1] = oldStat.CrawlerUpdatedAt
		oldMap[oldStat.ServerName] = info
	}
	var changed []domain.CrawlerStat
	for _, newStat := range newStats {
		oldValue, exists := oldMap[newStat.ServerName]
		if exists && newStat.CrawlerUpdatedAt.Equal(oldValue[1].(time.Time)) {
			continue
		}
		if !exists || oldValue[0] != newStat.Minutes {
			changed = append(changed, newStat)
		}
	}

	return changed
}
