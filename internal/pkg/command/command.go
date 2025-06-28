package command

import (
	"context"
	"fmt"
	"github.com/muerewa/hekaton/internal/pkg/helpers"
	"log"
	"os/exec"
	"strings"
	"time"
)

func RunBashCommand(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// Run bash command with retry and timeout logic
func VerifyBash(name, bash, timeout string, retries int) (string, error) {
	retries = max(1, retries) // Retry count
	var (
		result string
		err    error
	)
	tm, err := helpers.ParseDurationWithDefaults(timeout)
	if err != nil {
		return "", fmt.Errorf("%s: timeout format error: %v; exit", name, err)
	}
	for attempt := 0; attempt < retries; attempt++ {
		result, err = RunCommandWithTimeout(bash, time.Duration(tm))
		if err != nil {
			// Add info if there is extra retries
			retrySuffix := ""
			if attempt < retries-1 {
				retrySuffix = "; retrying..."
			} // add suffix
			log.Printf("%s: command error (attempt %d/%d): %v%s",
				name, attempt+1, retries, err, retrySuffix)
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
