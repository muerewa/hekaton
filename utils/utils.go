package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/muerewa/hekaton/structs"
)

func RunBashCommand(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// ParseDurationWithDefaults parses duration strings with intelligent defaults
func ParseDurationWithDefaults(s string) (time.Duration, error) {
	// Handle empty string - default to 1 second
	if s == "" {
		return time.Second, nil
	}

	// Try parsing as a standard duration first (e.g., "1m", "2h", "30s")
	if duration, err := time.ParseDuration(s); err == nil {
		return duration, nil
	}

	// If that fails, try parsing as a plain number (assume seconds)
	if seconds, err := strconv.Atoi(s); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	// If both fail, return the original parsing error
	_, err := time.ParseDuration(s)
	return 0, err
}

func Compare(result string, comp *structs.Compare) (bool, error) {
	// Integer comparison
	if numRes, err := strconv.Atoi(result); err == nil {
		if numComp, ok := comp.Value.(int); ok {
			switch comp.Operator {
			case ">":
				return numRes > numComp, nil
			case ">=":
				return numRes >= numComp, nil
			case "<":
				return numRes < numComp, nil
			case "<=":
				return numRes <= numComp, nil
			}
		}
	}

	// String comparison
	strComp := fmt.Sprintf("%v", comp.Value)
	switch comp.Operator {
	case "==":
		return result == strComp, nil
	case "!=":
		return result != strComp, nil
	default:
		return false, fmt.Errorf("unsupported operator")
	}
}
