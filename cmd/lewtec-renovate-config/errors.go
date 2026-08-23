package main

import (
	"errors"
	"log/slog"
)

// Sentinel errors for ReportError / wrapping with fmt.Errorf %w.
var ErrGhCLINotFound = errors.New("gh CLI not found. Please install it: https://cli.github.com/")
var ErrGhCLINotAuth = errors.New("gh CLI is not authenticated. Please run: gh auth login")
var ErrRepoProcess = errors.New("processing repository")
var ErrBatchProcessing = errors.New("batch repository processing")

// ReportError reports an error to the centralized error tracking system.
// currently uses slog, but is designed to be easily extensible to Sentry or others.
func ReportError(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	// Prepend error to args to ensure it's always logged
	allArgs := append([]any{"error", err}, args...)
	slog.Error(msg, allArgs...)
}
