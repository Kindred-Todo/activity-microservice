package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Kindred-Todo/activity-microservice/xslog"
)

func fatal(ctx context.Context, msg string, err error) {
	slog.LogAttrs(
		ctx,
		slog.LevelError,
		msg,
		xslog.Error(err),
	)
	os.Exit(1)
}