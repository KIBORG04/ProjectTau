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

// ---- ACCU (time-weighted average concurrent users) ----

// calcACCU computes the time-weighted average concurrent users for a given period.
//
// For each round overlapping [periodStart, periodEnd):
//
//	contribution = round.players * overlap_duration
//
// ACCU = total_player_seconds / period_duration_seconds
func calcACCU(rounds []Round, periodStart, periodEnd time.Time) int {
	totalPlayerSeconds := 0.0
	periodDuration := periodEnd.Sub(periodStart).Seconds()
	if periodDuration <= 0 {
		return 0
	}

	for _, r := range rounds {
		s, e, ok := clipRound(r, periodStart, periodEnd)
		if !ok {
			continue
		}
		totalPlayerSeconds += float64(r.Players) * e.Sub(s).Seconds()
	}

	return int(math.Round(totalPlayerSeconds / periodDuration))
}

// ---- PCCU (peak concurrent users, sweep-line) ----

type sweepEvent struct {
	time  time.Time
	delta int // +players at start, -players at end
	isEnd bool
}

// calcPCCU computes the peak concurrent users within [periodStart, periodEnd) using sweep-line.
func calcPCCU(rounds []Round, periodStart, periodEnd time.Time) int {
	var events []sweepEvent

	for _, r := range rounds {
		s, e, ok := clipRound(r, periodStart, periodEnd)
		if !ok {
			continue
		}
		events = append(events, sweepEvent{time: s, delta: r.Players, isEnd: false})
		events = append(events, sweepEvent{time: e, delta: -r.Players, isEnd: true})
	}

	if len(events) == 0 {
		return 0
	}

	// Sort by time; for equal times, ends (isEnd=true) come before starts
	// so that a round ending at T and a round starting at T don't overlap.
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
	current := 0
	for _, ev := range events {
		current += ev.delta
		if current > maxConcurrent {
			maxConcurrent = current
		}
	}

	return maxConcurrent
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
		accuMap[label] = calcACCU(rounds, weekStart, weekEnd)
		pccuMap[label] = calcPCCU(rounds, weekStart, weekEnd)

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
		accuMap[label] = calcACCU(rounds, dayStart, dayEnd)
		pccuMap[label] = calcPCCU(rounds, dayStart, dayEnd)
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

			accu := calcACCU(rounds, slotStart, slotEnd)
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
