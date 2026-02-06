package main

import (
	"log/slog"
)

/**
 * Reports an error to the centralized error tracking system.
 * currently uses slog, but is designed to be easily extensible to Sentry or others.
 *
 * @param err The error to report.
 * @param msg A descriptive message.
 * @param args Additional key-value pairs for context (slog style).
 */
func ReportError(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	// Prepend error to args to ensure it's always logged
	allArgs := append([]any{"error", err}, args...)
	slog.Error(msg, allArgs...)
}
