package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/muerewa/hekaton/structs"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) ([]structs.Monitor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var monitors []structs.Monitor
	err = yaml.Unmarshal(data, &monitors)
	return monitors, err
}

func RunBashCommand(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func Compare(result string, comp *structs.Compare) (bool, error) {
	// Попытка численного сравнения
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

	// Строковое сравнение
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

func VerifyBash(monitor *structs.Monitor) (string, error) {
	retries := max(1, monitor.Retries)
	timeout := time.Duration(max(1, monitor.Timeout)) * time.Second
	var (
		result string
		err    error
	)

	for attempt := 0; attempt < retries; attempt++ {
		result, err = RunCommandWithTimeout(monitor.Bash, time.Duration(timeout))
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

func RunCommandWithTimeout(command string, timeout time.Duration) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure resources are released

	// Prepare the command with context
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	// Execute and capture combined output
	output, err := cmd.CombinedOutput()

	// Handle context-specific errors
	if ctx.Err() == context.DeadlineExceeded {
		return strings.TrimSpace(string(output)), fmt.Errorf("command timed out after %v: %w", timeout, ctx.Err())
	}

	// Return other execution errors (exit status != 0)
	if err != nil {
		return strings.TrimSpace(string(output)), fmt.Errorf("command execution failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
