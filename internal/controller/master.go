package controller

import (
	"ssstatistics/internal/service"
	"time"
)

const currentStatistics = "02-27-2022"

func Run() {
	startDate, _ := time.Parse("01-02-2006", currentStatistics)

	collector := service.Collector{}
	collector.CollectUrls(startDate)
}
