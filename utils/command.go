package utils

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/muerewa/hekaton/structs"
)

// Run bash command with retry and timeout logic
func VerifyBash(monitor *structs.Monitor) (string, error) {
	retries := max(1, monitor.Retries) // Retry count
	var (
		result string
		err    error
	)
	timeout, err := ParseDurationWithDefaults(monitor.Timeout)
	if err != nil {
		return "", fmt.Errorf("%s: timeout format error: %v; exit", monitor.Name, err)
	}
	for attempt := 0; attempt < retries; attempt++ {
		result, err = RunCommandWithTimeout(monitor.Bash, time.Duration(timeout))
		if err != nil {
			// Add info if there is extra retries
			retrySuffix := ""
			if attempt < retries-1 {
				retrySuffix = "; retrying..."
			} // add suffix
			log.Printf("%s: command error (attempt %d/%d): %v%s",
				monitor.Name, attempt+1, retries, err, retrySuffix)
		}
	}
	return result, err
}

// Run command with timeout
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
