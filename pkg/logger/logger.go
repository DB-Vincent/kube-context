/*
 * kube-context
 *
 * Copyright (C) 2024 Vincent De Borger
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	debug bool
}

// New creates a new ErrorHandler instance
func New(debug bool) *Logger {
	return &Logger{debug: debug}
}

// Handle processes and displays a user-friendly error message
func (l *Logger) Handle(errType ErrorType, err error, args ...interface{}) {
	// Format message with any provided arguments
	message := errType.Message
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}


	// Print user-friendly message with appropriate prefix
	if errType.Level == Info {
		prefix := l.getPrefix(errType.Level)
		fmt.Fprintf(os.Stdout, "%s %s\n", prefix, message)
	} else {
		prefix := l.getPrefix(errType.Level)
		fmt.Fprintf(os.Stderr, "%s %s\n", prefix, message)
	}

	// Print debug information if enabled
	if l.IsDebug() && err != nil {
		fmt.Fprintf(os.Stderr, "Debug details: %v\n", err)
	}

	// Exit on fatal errors
	if errType.Level == Fatal {
		os.Exit(1)
	}
}

// getPrefix returns the appropriate prefix for the error level
func (l *Logger) getPrefix(level ErrorLevel) string {
	switch level {
	case Info:
		return "ℹ️"
	case Warning:
		return "⚠️"
	case Error:
		return "❌"
	case Fatal:
		return "💀"
	default:
		return "❓"
	}
}

// IsDebug returns true if debug mode is enabled
func (l *Logger) IsDebug() bool {
	return l.debug
}