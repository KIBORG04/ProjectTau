package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const maxRoundDuration = 4 * time.Hour

func main() {
	configPath := "config/online_stats.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := LoadOnlineStatsConfig(configPath)
	if err != nil {
		log.Printf("Warning: could not load config (%v), using defaults\n", err)
		cfg = &OnlineStatsConfig{
			OutputPath: "web/static/data/online_stats.json",
		}
	}

	log.Println("Loading rounds data...")
	rounds := loadRounds(cfg)

	log.Printf("Loaded %d rounds\n", len(rounds))
	if len(rounds) == 0 {
		log.Println("No rounds to process, exiting.")
		return
	}

	now := time.Now().UTC()

	log.Println("Calculating weeks data (ACCU/PCCU)...")
	weeksData := CalcWeeks(rounds)

	log.Println("Calculating last 90 days data (ACCU/PCCU)...")
	last90Data := CalcLast90Days(rounds, now)

	log.Println("Calculating daytime data (ACCU)...")
	daytimeData := CalcDaytime(rounds)

	output := OnlineStatsOutput{
		GeneratedAt: now.Format(time.RFC3339),
		Weeks:       weeksData,
		Last90Days:  last90Data,
		Daytime:     daytimeData,
	}

	if err := writeJSON(cfg.OutputPath, output); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	log.Printf("Successfully wrote online stats to %s\n", cfg.OutputPath)
}

func writeJSON(path string, data any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// loadRounds loads round data. It tries DB via SSH, then falls back to Mock data.
func loadRounds(cfg *OnlineStatsConfig) []Round {
	// Try loading existing state/cache
	state, err := LoadState(cfg.StatePath)
	if err != nil {
		log.Printf("Warning: failed to load state from %s: %v\n", cfg.StatePath, err)
		// Continue with empty state
		state = &StatsState{}
	}

	// If SSHHost is provided, fetch new rounds from DB
	if cfg.RemoteDB.SSHHost != "" {
		log.Printf("State loaded. Cached rounds: %d, Last parsed ID: %d\n", len(state.Rounds), state.LastRoundID)

		newRounds, err := FetchNewRoundsFromDB(cfg, state.LastRoundID)
		if err != nil {
			log.Printf("Error fetching rounds from DB: %v\n", err)
			log.Println("Proceeding with cached rounds only.")
		} else if len(newRounds) > 0 {
			state.Rounds = append(state.Rounds, newRounds...)

			// Find new max ID
			maxID := state.LastRoundID
			for _, r := range newRounds {
				if r.ID > maxID {
					maxID = r.ID
				}
			}
			state.LastRoundID = maxID

			if err := SaveState(cfg.StatePath, state); err != nil {
				log.Printf("Warning: failed to save state: %v\n", err)
			} else {
				log.Printf("Saved updated state with %d total rounds, last ID: %d\n", len(state.Rounds), state.LastRoundID)
			}
		}

		filtered := filterInvalidRounds(state.Rounds)
		return filtered
	}

	if len(state.Rounds) > 0 {
		filtered := filterInvalidRounds(state.Rounds)
		return filtered
	}

	log.Println("Using mock data (DB connection not configured)")
	return generateMockRounds()
}

// filterInvalidRounds removes invalid rounds and caps duration at maxRoundDuration.
func filterInvalidRounds(rounds []Round) []Round {
	zeroTime := time.Time{}
	var filtered []Round
	var skipped int

	for _, r := range rounds {
		// Skip rounds with zero/invalid dates
		if r.StartDatetime.Equal(zeroTime) || r.EndDatetime.Equal(zeroTime) {
			skipped++
			continue
		}

		// Skip rounds with zero port
		if r.ServerPort == 0 {
			skipped++
			continue
		}

		// Skip rounds with non-positive duration
		if !r.StartDatetime.Before(r.EndDatetime) {
			skipped++
			continue
		}

		// Cap duration at maxRoundDuration (take last N hours) to handle hung servers
		if r.EndDatetime.Sub(r.StartDatetime) > maxRoundDuration {
			r.StartDatetime = r.EndDatetime.Add(-maxRoundDuration)
		}

		filtered = append(filtered, r)
	}

	if skipped > 0 {
		log.Printf("Filtered out %d invalid rounds (zero dates, zero ports, negative duration)\n", skipped)
	}

	return filtered
}

// generateMockRounds creates realistic mock data for testing.
// Simulates rounds across 3 servers over ~100 days with typical daily patterns.
func generateMockRounds() []Round {
	var rounds []Round
	id := 1

	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -100)

	ports := []int{1337, 1338, 1339}

	for d := startDate; d.Before(now); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()

		for _, port := range ports {
			// Number of rounds per day per server (more on weekends)
			numRounds := 8
			if weekday == time.Saturday || weekday == time.Sunday {
				numRounds = 12
			}

			// Simulate activity from ~8:00 to ~3:00 next day (UTC)
			roundStart := d.Add(8 * time.Hour) // first round at 8:00

			for r := 0; r < numRounds; r++ {
				hour := roundStart.Hour()

				// Players depend on time of day (peak at 18-22)
				var players int
				switch {
				case hour >= 18 && hour <= 22:
					players = 30 + (id % 20) // 30-49 during peak
				case hour >= 14 && hour < 18:
					players = 15 + (id % 15) // 15-29 afternoon
				case hour >= 10 && hour < 14:
					players = 8 + (id % 10) // 8-17 morning
				default:
					players = 2 + (id % 5) // 2-6 night
				}

				// Weekend bonus
				if weekday == time.Saturday || weekday == time.Sunday {
					players = int(float64(players) * 1.3)
				}

				// Round duration: 40-120 min
				duration := time.Duration(40+(id%80)) * time.Minute

				roundEnd := roundStart.Add(duration)

				rounds = append(rounds, Round{
					ID:            id,
					StartDatetime: roundStart,
					EndDatetime:   roundEnd,
					ServerPort:    port,
					Players:       players,
				})

				id++

				// Next round starts 5-20 min after this one ends (lobby time)
				gap := time.Duration(5+(id%15)) * time.Minute
				roundStart = roundEnd.Add(gap)
			}
		}
	}

	return rounds
}
