package controller

import (
	"ssstatistics/internal/repository"
	"ssstatistics/internal/service/cleaning"
	"ssstatistics/internal/service/maps"
	"ssstatistics/internal/service/parser"
	"ssstatistics/internal/service/stats/ckey_statistics/crawler"
	"ssstatistics/internal/service/stats/ckey_statistics/mmr"
	"ssstatistics/internal/service/tops"
)

type GeneralUpdater func() []string

var GeneralUpdaters []GeneralUpdater

func InitializeGeneralUpdaters() {
	GeneralUpdaters = append(GeneralUpdaters, parser.RunRoundCollector)
}

func StartRegularUpdaters() []string {
	var logs []string

	for _, callback := range GeneralUpdaters {
		callbackLogs := callback()
		logs = append(logs, callbackLogs...)
	}
	return logs
}

var DBUpdaters []DBUpdater

type DBUpdater func() []string

func InitializeDBUpdaters() {
	DBUpdaters = append(DBUpdaters, tops.ParseTopData)
	DBUpdaters = append(DBUpdaters, mmr.ParseMMR)
	DBUpdaters = append(DBUpdaters, maps.FixMaxShit)
	DBUpdaters = append(DBUpdaters, cleaning.CleanAnnounces)
	DBUpdaters = append(DBUpdaters, cleaning.CleanDuplicatedFlavors)
	DBUpdaters = append(DBUpdaters, repository.RefreshMaterializedViews)
	DBUpdaters = append(DBUpdaters, crawler.SecretlyUpdateSomePlayers)
}

func StartDBUpdaters() []string {
	var logs []string

	for _, callback := range DBUpdaters {
		callbackLogs := callback()
		logs = append(logs, callbackLogs...)
	}
	return logs
}

func StartUpdaters() []string {
	var logs []string
	logs = append(logs, StartRegularUpdaters()...)
	logs = append(logs, StartDBUpdaters()...)
	return logs
}

func InitializeUpdaters() {
	InitializeGeneralUpdaters()
	InitializeDBUpdaters()
}
