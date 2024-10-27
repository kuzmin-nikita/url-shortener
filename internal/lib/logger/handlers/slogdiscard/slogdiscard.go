package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type discardHandler struct{}

func NewDiscardHandler() *discardHandler {
	return &discardHandler{}
}

func (h *discardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (h *discardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h *discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *discardHandler) WithGroup(name string) slog.Handler {
	return h
}
