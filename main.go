package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/muerewa/hekaton/core"
	"github.com/muerewa/hekaton/utils"
)

func main() {
	monitors, err := utils.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for _, monitor := range monitors {
		go core.RunMonitor(ctx, &monitor)
	}

	<-sigChan
	log.Println("Получен сигнал остановки, завершаем работу...")
	cancel()
}
