package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wzshiming/gh-gpt/pkg/cmd"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		cancel()
		<-sigChan
		os.Exit(1)
	}()

	command := cmd.NewCommand()
	err := command.ExecuteContext(ctx)
	if err != nil {
		slog.Error("failed to execute command", "err", err)
		os.Exit(1)
	}
}
