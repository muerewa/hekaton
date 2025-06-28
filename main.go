package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/muerewa/hekaton/core"
	"github.com/muerewa/hekaton/utils"
)

func main() {

	configPath := flag.String("config", "config.yaml", "Path of config file")
	flag.Parse()

	// Parsing config
	monitors, err := utils.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Creating config for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signals handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// For every rule run a goroutine
	for _, monitor := range monitors {
		go core.RunMonitor(ctx, &monitor)
	}

	<-sigChan // Wait for a stop signal
	log.Println("Получен сигнал остановки, завершаем работу...")
	cancel()
}
