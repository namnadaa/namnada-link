package main

import (
	"URLbot/pkg/clients/telegram"
	eventconsumer "URLbot/pkg/consumer/event-consumer"
	tgEvents "URLbot/pkg/events/telegram"
	"URLbot/pkg/storage/memory"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tgClient := telegram.NewClient(mustParseFlags())

	storage := memory.New()

	eventProcessor := tgEvents.New(tgClient, storage)

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		batchSize = 100
		slog.Warn("Invalid or missing BATCH_SIZE, using default", "default", batchSize)
	}

	slog.Info("Service started")

	consumer := eventconsumer.New(eventProcessor, eventProcessor, batchSize)

	err = consumer.Start(ctx)
	slog.Info("Service is shutting down")

	if err != nil {
		slog.Error("Service stopped with error", "err", err)
		os.Exit(1)
	} else {
		slog.Info("Service stopped gracefully")
	}
}

// mustParseFlags parses required command-line flags (scheme, host, token).
func mustParseFlags() (string, string, string) {
	scheme := flag.String("tg-bot-scheme", "", "Scheme for Telegram Bot API (e.g., https)")
	host := flag.String("tg-bot-host", "", "Telegram Bot API host (e.g., api.telegram.org)")
	token := flag.String("tg-bot-token", "", "Access token for Telegram bot")

	flag.Parse()

	if *scheme == "" || *host == "" || *token == "" {
		slog.Error("Missing required flags", "scheme", *scheme, "host", *host, "token_present", *token != "")
		os.Exit(1)
	}

	return *scheme, *host, *token
}
