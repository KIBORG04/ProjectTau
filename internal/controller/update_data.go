package controller

import (
	"ssstatistics/internal/service/cleaning"
	"ssstatistics/internal/service/mmr"
	"ssstatistics/internal/service/parser"
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
		for _, s := range callbackLogs {
			logs = append(logs, s)
		}
	}
	return logs
}

var DBUpdaters []DBUpdater

// DBUpdater
//  TODO: Возможно, стоит сделать поттягивание всей БД и передача её в каждую из функций,
//  TODO: чтобы они не делали свои запросы по кд
type DBUpdater func() []string

func InitializeDBUpdaters() {
	DBUpdaters = append(DBUpdaters, tops.ParseTopData)
	DBUpdaters = append(DBUpdaters, mmr.ParseMMR)
	DBUpdaters = append(DBUpdaters, cleaning.CleanAnnounces)
}

func StartDBUpdaters() []string {
	var logs []string

	for _, callback := range DBUpdaters {
		callbackLogs := callback()
		for _, s := range callbackLogs {
			logs = append(logs, s)
		}
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