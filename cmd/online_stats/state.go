package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StatsState struct {
	LastRoundID int     `json:"last_round_id"`
	Rounds      []Round `json:"rounds"`
}

func LoadState(path string) (*StatsState, error) {
	if path == "" {
		return &StatsState{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &StatsState{}, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state StatsState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	return &state, nil
}

func SaveState(path string, state *StatsState) error {
	if path == "" {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize state: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}
