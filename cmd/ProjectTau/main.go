package main

import (
	"ssstatistics/internal/config"
	"ssstatistics/internal/controller"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

func main() {
	config.LoadConfigurations()

	stats.PopulatePositions()

	r.CreateConnection()
	controller.Run()
}
