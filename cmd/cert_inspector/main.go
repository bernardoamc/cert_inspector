package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bernardoamc/cert-inspector/internal/logger"
	"github.com/bernardoamc/cert-inspector/internal/runner"
)

func runWithRestart(binary string, args []string, delay int, stdoutLogger, stderrLogger *logger.DailyLogger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		fmt.Println("Starting gungnir process...")
		ctx, cancel := context.WithCancel(context.Background())
		errCh := make(chan error, 1)

		go func() {
			errCh <- runner.RunProcess(ctx, binary, args, stdoutLogger, stderrLogger)
		}()

		select {
		case sig := <-stop:
			fmt.Printf("\nReceived signal: %v. Stopping...\n", sig)
			cancel()
			time.Sleep(1 * time.Second)
			return
		case err := <-errCh:
			if err != nil {
				fmt.Printf("Process error: %v\n", err)
			} else {
				fmt.Println("Process exited normally")
			}
			fmt.Printf("Restarting in %d seconds...\n", delay)
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}

func main() {
	binary := flag.String("binary", "./gungnir", "Path to gungnir binary")
	targets := flag.String("targets", "targets.txt", "Path to targets file")
	logDir := flag.String("log-dir", "logs", "Directory for logs")
	stdoutBase := flag.String("stdout-output", "domains", "Base name for stdout logs")
	stderrBase := flag.String("stderr-output", "errors", "Base name for stderr logs")
	restartDelay := flag.Int("restart-delay", 2, "Seconds before restart")
	flag.Parse()

	args := []string{"-r", *targets}

	stdoutLogger, err := logger.NewDailyLogger(*logDir, *stdoutBase)
	if err != nil {
		fmt.Printf("Error initializing stdout logger: %v\n", err)
		os.Exit(1)
	}
	defer stdoutLogger.Close()

	stderrLogger, err := logger.NewDailyLogger(*logDir, *stderrBase)
	if err != nil {
		fmt.Printf("Error initializing stderr logger: %v\n", err)
		os.Exit(1)
	}
	defer stderrLogger.Close()

	runWithRestart(*binary, args, *restartDelay, stdoutLogger, stderrLogger)
}
