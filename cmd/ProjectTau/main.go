package main

import (
	"ssstatistics/internal/config"
	"ssstatistics/internal/controller"
	r "ssstatistics/internal/repository"
)

func main() {
	config.LoadConfigurations()

	r.CreateConnection()
	controller.Run()
//  NEVERMIND
// 	var a domain.CultInfo
// 	r.Database.Preload(clause.Associations).Find(&a)
// 	fmt.Println(a)
}
