package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muerewa/hekaton/structs"
	"github.com/muerewa/hekaton/utils"
)

func VerifyBash(monitor *structs.Monitor) (string, error) {
	retries := max(1, monitor.Retries)
	var (
		result string
		err    error
	)

	for attempt := 0; attempt < retries; attempt++ {
		result, err = utils.RunBashCommand(monitor.Bash)
		if err != nil {
			// Добавляем информацию о попытке только если будут еще ретраи
			retrySuffix := ""
			if attempt < retries-1 {
				retrySuffix = "; retrying..."
			}
			log.Printf("%s: command error (attempt %d/%d): %v%s",
				monitor.Name, attempt+1, retries, err, retrySuffix)
		}
	}
	return result, err
}

func RunMonitor(ctx context.Context, monitor *structs.Monitor) {
	interval := time.Second
	if monitor.Interval > 0 {
		interval = time.Duration(monitor.Interval) * time.Second
	}
	ticker := time.Tick(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			result, err := VerifyBash(monitor)
			if err != nil {
				log.Printf("%s: command error: %v; exit...", monitor.Name, err)
				return
			}

			matched, err := utils.Compare(result, &monitor.Compare)
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
					}
				}
			}
		}
	}
}
