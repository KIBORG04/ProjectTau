package main

import (
	"ssstatistics/internal/config"
	"ssstatistics/internal/controller"
	db "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

func main() {
	config.LoadConfigurations()

	stats.PopulatePositions()

	db.CreateConnection()

	controller.InitializeRegularCallbacks()

	controller.Run()
}
