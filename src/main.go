package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ssstatistics/domain"
)

func main() {
	resp, _ := http.Get("https://stat.taucetistation.org/html/2022/03/02/round-49364/stat.json")

	var root domain.Root
	json.NewDecoder(resp.Body).Decode(&root)

	fmt.Print(root.Factions[0].CultInfo.RitenameByCount)
}
