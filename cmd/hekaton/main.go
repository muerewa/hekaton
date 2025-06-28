package main

import (
	"context"
	"flag"
	"github.com/muerewa/hekaton/internal/pkg/config"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/muerewa/hekaton/internal/app"
)

func main() {

	configPath := flag.String("config", "config.yaml", "Path of config file")
	flag.Parse()

	// Parsing config
	monitors, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Creating config for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signals handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// For every rule run a goroutine
	for _, monitor := range monitors {
		go app.RunMonitor(ctx, &app.MonitorRule{monitor, logger})
	}

	<-sigChan // Wait for a stop signal
	log.Println("Получен сигнал остановки, завершаем работу...")
	cancel()
}
