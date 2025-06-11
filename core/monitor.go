package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muerewa/hekaton/actions"
	"github.com/muerewa/hekaton/structs"
	"github.com/muerewa/hekaton/utils"
)

// Executes actions
func ExecuteActions(ctx context.Context, monitor *structs.Monitor, result string) {
	for _, action := range monitor.Actions {
		switch action.Type {
		case "bash":
			res, _ := utils.RunBashCommand(action.Params["command"])
			fmt.Println(res)
		case "telegram":
			// Send message to Telegram
			err := actions.SendTelegramMessage(monitor.Name, action.Params, result)
			if err != nil {
				log.Printf("%s: telegram error: %v", monitor.Name, err)
			}
		case "email":
			err := actions.SendEmail(action.Params, result)
			if err != nil {
				log.Printf("%s: email error: %v", monitor.Name, err)
			}
		}
	}
}

// Processes one monitor tick
func ProcessMonitorTick(ctx context.Context, monitor *structs.Monitor) error {
	result, err := utils.VerifyBash(monitor) // Run bash command
	if err != nil {
		return fmt.Errorf("%s: command error: %v; exit", monitor.Name, err)
	}

	matched, err := utils.Compare(result, &monitor.Compare) // Compare result
	if err != nil {
		return fmt.Errorf("%s: compare error: %v", monitor.Name, err)
	}

	if matched {
		ExecuteActions(ctx, monitor, result)
	}
	return nil
}

// Main monitor gouroutine function
func RunMonitor(ctx context.Context, monitor *structs.Monitor) {
	// Setting interval value
	interval, err := utils.ParseDurationWithDefaults(monitor.Interval)
	if err != nil {
		log.Printf("%s: interval format error: %v; exit...", monitor.Name, err)
		return
	}
	// Create a ticker
	ticker := time.Tick(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			err = ProcessMonitorTick(ctx, monitor)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}
