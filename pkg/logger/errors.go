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

// ErrorLevel defines the severity of the error
type ErrorLevel int

const (
	Info ErrorLevel = iota
	Warning
	Error
	Fatal
)

// ErrorType represents common error scenarios
type ErrorType struct {
	Level   ErrorLevel
	Message string
}

// Common error types
var (
	ErrFatal = ErrorType{
		Level:   Fatal,
		Message: "This is a fatal error",
	}
	ErrError = ErrorType{
		Level:   Error,
		Message: "This is an error",
	}
	ErrWarning = ErrorType{
		Level:   Warning,
		Message: "This is a warning",
	}
	ErrInfo = ErrorType{
		Level:   Info,
		Message: "This is an info message",
	}
)