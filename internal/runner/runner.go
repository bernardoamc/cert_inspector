package runner

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/bernardoamc/cert-inspector/internal/logger"
)

func streamOutput(scanner *bufio.Scanner, tag string, logger *logger.DailyLogger) {
	for scanner.Scan() {
		line := scanner.Text()
		if tag != "" {
			line = tag + " " + line
		}
		fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), line)
		_ = logger.WriteLine(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] Scanner error: %v\n", tag, err)
	}
}

func RunProcess(ctx context.Context, binary string, args []string, stdoutLogger, stderrLogger *logger.DailyLogger) error {
	cmd := exec.CommandContext(ctx, binary, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		streamOutput(bufio.NewScanner(stdout), "", stdoutLogger)
	}()

	go func() {
		defer wg.Done()
		streamOutput(bufio.NewScanner(stderr), "[STDERR]", stderrLogger)
	}()

	err = cmd.Wait()
	wg.Wait()

	if err != nil {
		return fmt.Errorf("process exited: %w", err)
	}
	return nil
}
