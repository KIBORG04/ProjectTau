package main

import "time"

// Round represents a single game round from the remote database.
type Round struct {
	ID            int       `json:"id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	ServerPort    int       `json:"server_port"`
	Players       int       `json:"players"`
}

// OnlineStatsOutput is the top-level JSON structure written to the output file.
type OnlineStatsOutput struct {
	GeneratedAt string          `json:"generated_at"`
	Weeks       WeeksData       `json:"weeks"`
	Last90Days  Last90DaysData  `json:"last_90_days"`
	Daytime     DaytimeData     `json:"daytime"`
}

// WeeksData holds ACCU and PCCU per ISO week label (e.g., "2025-12").
type WeeksData struct {
	ACCU map[string]int `json:"accu"` // label → average concurrent
	PCCU map[string]int `json:"pccu"` // label → peak concurrent
}

// Last90DaysData holds ACCU and PCCU per day (e.g., "2025-03-01").
type Last90DaysData struct {
	ACCU map[string]int `json:"accu"`
	PCCU map[string]int `json:"pccu"`
}

// DaytimeData holds average concurrent players per 2-hour interval (0,2,4,...,22).
type DaytimeData struct {
	ACCU map[int]int `json:"accu"`
}
