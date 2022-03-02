package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ssstatistics/config"
	"ssstatistics/domain"
	r "ssstatistics/repository"

	"gorm.io/gorm/clause"
)

func main() {
	config.LoadConfigurations()

	r.CreateConnection()

	resp, err := http.Get("https://stat.taucetistation.org/html/2022/03/02/round-49364/stat.json")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}

	var root domain.Root
	json.NewDecoder(resp.Body).Decode(&root)
	r.Database.Save(&root)

	var a domain.CultInfo
	r.Database.Preload(clause.Associations).Find(&a)
	fmt.Println(a)
}
