package core

import (
	"fmt"
	"log"
	"time"

	"github.com/muerewa/hekaton/structs"
	"github.com/muerewa/hekaton/utils"
)

func RunMonitor(monitor structs.Monitor) {
	interval := time.Second
	if monitor.Interval > 0 {
		interval = time.Duration(monitor.Interval) * time.Second
	}

	for range time.Tick(interval) {
		result, err := utils.RunBashCommand(monitor.Bash)
		if err != nil {
			log.Printf("%s: command error: %v", monitor.Name, err)
			continue
		}

		matched, err := utils.Compare(result, monitor.Compare)
		if err != nil {
			log.Printf("%s: compare error: %v", monitor.Name, err)
			continue
		}

		if matched {
			for _, action := range monitor.Actions {
				switch action.Type {
				case "bash":
					fmt.Println(utils.RunBashCommand(action.Params["command"]))
				}
			}
		}
	}
}
