package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Kindred-Todo/activity-microservice/xslog"
)

func fatal(ctx context.Context, msg string, err error) {
	fmt.Println(msg, err)
	slog.LogAttrs(
		ctx,
		slog.LevelError,
		msg,
		xslog.Error(err),
	)
	os.Exit(1)
}