package chartjs

import (
	"encoding/json"
	"fmt"
	"html/template"
)

type Config struct {
	Type    string        `json:"type"`
	Data    DataStructure `json:"data,omitempty"`
	Options string        `json:"options,omitempty"`
}

type DataStructure struct {
	Labels   []string   `json:"labels"`
	Datasets []*Dataset `json:"datasets"`
}

func New(_type string) *Config {
	config := Config{
		Type: _type,
	}

	return &config
}

func (c *Config) SetLabels(labels []string) *Config {
	c.Data.Labels = labels
	return c
}

func (c *Config) AddDataset(dataset *Dataset) *Config {
	c.Data.Datasets = append(c.Data.Datasets, dataset)
	return c
}

func (c *Config) String() template.JS {
	str, err := json.Marshal(&c)
	if err != nil {
		fmt.Println(err)
	}
	return template.JS(str)
}
