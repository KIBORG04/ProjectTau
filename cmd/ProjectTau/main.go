package main

import (
	"ssstatistics/internal/bots/telegram"
	"ssstatistics/internal/config"
	"ssstatistics/internal/controller"
	db "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

func main() {
	config.LoadConfigurations()

	stats.PopulatePositions()

	db.CreateConnection()

	telegram.Initialize()

	controller.InitializeUpdaters()

	controller.Run()
}
