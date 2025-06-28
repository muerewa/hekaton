package main

import (
	"context"
	"flag"
	"github.com/muerewa/hekaton/internal/pkg/config"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/muerewa/hekaton/internal/app"
	_ "net/http/pprof"
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

	var wg sync.WaitGroup

	// For every rule run a goroutine
	for _, monitor := range monitors {
		wg.Add(1)
		go app.RunMonitor(ctx, &wg, &app.MonitorRule{
			Monitor: monitor, Log: logger})
	}
	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()
	select {
	case <-sigChan:
		logger.Info("Получен сигнал остановки, завершаем работу...")
		cancel()
		<-waitChan
	case <-waitChan:
		logger.Info("Все мониторы были завершены")
	}
}
