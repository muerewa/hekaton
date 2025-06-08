package main

import (
	"log"

	"github.com/muerewa/hekaton/core"
	"github.com/muerewa/hekaton/utils"
)

func main() {
	monitors, err := utils.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	for _, monitor := range monitors {
		go core.RunMonitor(monitor)
	}

	// Бесконечное ожидание
	select {}
}
