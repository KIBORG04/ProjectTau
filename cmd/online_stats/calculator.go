package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// ---- Helpers ----

func isoYearWeek(t time.Time) (int, int) {
	year, week := t.ISOWeek()
	return year, week
}

func weekLabel(year, week int) string {
	return fmt.Sprintf("%d-%02d", year, week)
}

func dayLabel(t time.Time) string {
	return t.Format("2006-01-02")
}

// startOfISOWeek returns the Monday 00:00:00 UTC of the given ISO week.
func startOfISOWeek(year, week int) time.Time {
	// Jan 4 is always in ISO week 1
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	dow := jan4.Weekday()
	if dow == 0 {
		dow = 7
	}
	monday := jan4.AddDate(0, 0, -int(dow)+1) // Monday of week 1
	return monday.AddDate(0, 0, (week-1)*7)
}

// ---- Clipping ----

// clipRound clips a round to the interval [start, end). Returns clamped start/end and false if no overlap.
func clipRound(r Round, start, end time.Time) (time.Time, time.Time, bool) {
	s := r.StartDatetime
	e := r.EndDatetime
	if s.Before(start) {
		s = start
	}
	if e.After(end) {
		e = end
	}
	if !s.Before(e) {
		return s, e, false
	}
	return s, e, true
}

type sweepEvent struct {
	time       time.Time
	port       int
	roundIndex int
	players    int
	isEnd      bool
}

// calcStats computes both ACCU and PCCU for a given period using a server-grouped sweep-line.
// To prevent duplicated/hung rounds on the same server from falsely adding up,
// the online players for a given server port is the MAX of its active rounds.
// Global online is the SUM of these per-port maximums.
// ACCU is the time-weighted average (Area Under Curve) of this global online.
func calcStats(rounds []Round, periodStart, periodEnd time.Time) (int, int) {
	var events []sweepEvent

	for i, r := range rounds {
		s, e, ok := clipRound(r, periodStart, periodEnd)
		if !ok {
			continue
		}
		events = append(events, sweepEvent{time: s, port: r.ServerPort, roundIndex: i, players: r.Players, isEnd: false})
		events = append(events, sweepEvent{time: e, port: r.ServerPort, roundIndex: i, players: r.Players, isEnd: true})
	}

	periodDuration := periodEnd.Sub(periodStart).Seconds()
	if len(events) == 0 || periodDuration <= 0 {
		return 0, 0
	}

	// Sort by time; for equal times, ends (isEnd=true) come before starts
	sort.Slice(events, func(i, j int) bool {
		if events[i].time.Equal(events[j].time) {
			if events[i].isEnd != events[j].isEnd {
				return events[i].isEnd
			}
			return false
		}
		return events[i].time.Before(events[j].time)
	})

	maxConcurrent := 0
	totalPlayerSeconds := 0.0
	activeByPort := make(map[int]map[int]int)

	lastTime := periodStart
	current := 0 // Global current online

	for i := 0; i < len(events); {
		t := events[i].time

		// Add area for the interval [lastTime, t]
		if current > 0 && t.After(lastTime) {
			totalPlayerSeconds += float64(current) * t.Sub(lastTime).Seconds()
		}

		// Process all events occurring exactly at time `t`
		for i < len(events) && events[i].time.Equal(t) {
			ev := events[i]
			if activeByPort[ev.port] == nil {
				activeByPort[ev.port] = make(map[int]int)
			}
			if ev.isEnd {
				delete(activeByPort[ev.port], ev.roundIndex)
			} else {
				activeByPort[ev.port][ev.roundIndex] = ev.players
			}
			i++
		}

		// Recalculate global current from the active ports
		current = 0
		for _, roundsOnPort := range activeByPort {
			portMax := 0
			for _, players := range roundsOnPort {
				if players > portMax {
					portMax = players
				}
			}
			current += portMax
		}

		if current > maxConcurrent {
			maxConcurrent = current
		}

		lastTime = t
	}

	// Add any remaining area from the last event until periodEnd
	if current > 0 && lastTime.Before(periodEnd) {
		totalPlayerSeconds += float64(current) * periodEnd.Sub(lastTime).Seconds()
	}

	accu := int(math.Round(totalPlayerSeconds / periodDuration))
	return accu, maxConcurrent
}

// ---- Chart Calculations ----

// CalcWeeks calculates ACCU and PCCU per ISO week across all rounds.
func CalcWeeks(rounds []Round) WeeksData {
	if len(rounds) == 0 {
		return WeeksData{ACCU: map[string]int{}, PCCU: map[string]int{}}
	}

	// Find the global min/max dates
	minTime := rounds[0].StartDatetime
	maxTime := rounds[0].EndDatetime
	for _, r := range rounds {
		if r.StartDatetime.Before(minTime) {
			minTime = r.StartDatetime
		}
		if r.EndDatetime.After(maxTime) {
			maxTime = r.EndDatetime
		}
	}

	accuMap := make(map[string]int)
	pccuMap := make(map[string]int)

	// Iterate over every ISO week from minTime to maxTime
	y, w := isoYearWeek(minTime)
	for {
		weekStart := startOfISOWeek(y, w)
		weekEnd := weekStart.AddDate(0, 0, 7)

		if weekStart.After(maxTime) {
			break
		}

		label := weekLabel(y, w)
		accu, pccu := calcStats(rounds, weekStart, weekEnd)
		accuMap[label] = accu
		pccuMap[label] = pccu

		// Advance to next week
		next := weekStart.AddDate(0, 0, 7)
		y, w = isoYearWeek(next)
	}

	return WeeksData{ACCU: accuMap, PCCU: pccuMap}
}

// CalcLast90Days calculates ACCU and PCCU per day for the last 90 days.
func CalcLast90Days(rounds []Round, now time.Time) Last90DaysData {
	accuMap := make(map[string]int)
	pccuMap := make(map[string]int)

	// "now" is today at midnight UTC; we exclude today and go back 90 days
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := today           // exclusive: today is not included
	startDate := today.AddDate(0, 0, -90)

	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dayStart := d
		dayEnd := d.AddDate(0, 0, 1)
		label := dayLabel(d)
		accu, pccu := calcStats(rounds, dayStart, dayEnd)
		accuMap[label] = accu
		pccuMap[label] = pccu
	}

	return Last90DaysData{ACCU: accuMap, PCCU: pccuMap}
}

// CalcDaytime calculates the average concurrent players per 2-hour interval
// over all days that have data.
//
// For each day, compute ACCU for each 2h slot [0:00-2:00), [2:00-4:00), ...
// Then average across all days.
func CalcDaytime(rounds []Round) DaytimeData {
	if len(rounds) == 0 {
		return DaytimeData{ACCU: map[int]int{}}
	}

	// Find global date range
	minTime := rounds[0].StartDatetime
	maxTime := rounds[0].EndDatetime
	for _, r := range rounds {
		if r.StartDatetime.Before(minTime) {
			minTime = r.StartDatetime
		}
		if r.EndDatetime.After(maxTime) {
			maxTime = r.EndDatetime
		}
	}

	startDay := time.Date(minTime.Year(), minTime.Month(), minTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(maxTime.Year(), maxTime.Month(), maxTime.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

	totalDays := int(endDay.Sub(startDay).Hours() / 24)
	if totalDays <= 0 {
		totalDays = 1
	}

	// Accumulate ACCU totals per 2h slot across all days
	slotSums := make(map[int]float64)  // slot → sum of ACCU values
	slotCounts := make(map[int]int)     // slot → number of days with data

	for d := startDay; d.Before(endDay); d = d.AddDate(0, 0, 1) {
		for slot := 0; slot < 24; slot += 2 {
			slotStart := d.Add(time.Duration(slot) * time.Hour)
			slotEnd := d.Add(time.Duration(slot+2) * time.Hour)

			accu, _ := calcStats(rounds, slotStart, slotEnd)
			slotSums[slot] += float64(accu)
			slotCounts[slot]++
		}
	}

	result := make(map[int]int)
	for slot := 0; slot < 24; slot += 2 {
		if slotCounts[slot] > 0 {
			result[slot] = int(math.Round(slotSums[slot] / float64(slotCounts[slot])))
		}
	}

	return DaytimeData{ACCU: result}
}
