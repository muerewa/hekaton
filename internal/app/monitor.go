package app

import (
	"context"
	"fmt"
	"github.com/muerewa/hekaton/internal/app/actions"
	"github.com/muerewa/hekaton/internal/pkg/command"
	"github.com/muerewa/hekaton/internal/pkg/helpers"
	"log"
	"log/slog"
	"time"
)

type Action struct {
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}

type Compare struct {
	Operator string      `yaml:"operator"`
	Value    interface{} `yaml:"value"` // Can be either string of int
}

type Monitor struct {
	Name     string   `yaml:"name"`
	Bash     string   `yaml:"bash"` // Bash command
	Compare  Compare  `yaml:"compare"`
	Actions  []Action `yaml:"actions"`
	Interval string   `yaml:"interval,omitempty"` // Interval: format - "1s", 2, "4m" etc
	Timeout  string   `yaml:"timeout,omitempty"`  // Timeout: format - "1s", 2, "4m" etc
	Retries  int      `yaml:"retries,omitempty"`  // Amount of retries
}

type MonitorRule struct {
	Monitor
	log *slog.Logger
}

// Executes actions
func ExecuteActions(ctx context.Context, monitor *Monitor, result string) {
	for _, action := range monitor.Actions {
		switch action.Type {
		case "bash":
			res, _ := command.RunBashCommand(action.Params["command"])
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
func ProcessMonitorTick(ctx context.Context, monitor *Monitor) error {
	result, err := command.VerifyBash(monitor.Name, monitor.Bash, monitor.Timeout, monitor.Retries) // Run bash command
	if err != nil {
		return fmt.Errorf("%s: command error: %v; exit", monitor.Name, err)
	}

	matched, err := helpers.CompareOperator(result, monitor.Compare.Operator, monitor.Compare.Value) // Compare result
	if err != nil {
		return fmt.Errorf("%s: compare error: %v", monitor.Name, err)
	}

	if matched {
		ExecuteActions(ctx, monitor, result)
	}
	return nil
}

// Main monitor gouroutine function
func RunMonitor(ctx context.Context, monitor *Monitor) {
	// Setting interval value
	interval, err := helpers.ParseDurationWithDefaults(monitor.Interval)
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
