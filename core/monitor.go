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

func RunMonitor(ctx context.Context, monitor *structs.Monitor) {
	// Setting interval value
	interval := time.Second
	if monitor.Interval > 0 {
		interval = time.Duration(monitor.Interval) * time.Second
	}
	// Create a ticker
	ticker := time.Tick(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			result, err := utils.VerifyBash(monitor) // Run bash command
			if err != nil {
				log.Printf("%s: command error: %v; exit...", monitor.Name, err)
				return
			}

			matched, err := utils.Compare(result, &monitor.Compare) // Compare result
			if err != nil {
				log.Printf("%s: compare error: %v", monitor.Name, err)
				continue
			}

			if matched {
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
					}
				}
			}
		}
	}
}
