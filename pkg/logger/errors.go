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
	ErrGetResource = ErrorType{
		Level: Fatal,
		Message: "Could not get %s resource(s) from Kubernetes cluster. Please verify that you have the correct permissions.",
	}
	ErrInitKubeconfig = ErrorType{
		Level:   Error,
		Message: "Failed to initialize kubeconfig",
	}
	ErrWriteKubeconfig = ErrorType{
		Level:   Error,
		Message: "Failed to write to kubeconfig",
	}
	ErrContextNotFound = ErrorType{
		Level:   Error,
		Message: "Could not find context in kubeconfig file! Found the following contexts: %q",
	}
	ErrSelectContext = ErrorType{
		Level:   Error,
		Message: "Error selecting context",
	}
	ErrUserInterrupt = ErrorType{
		Level:   Info,
		Message: "Alright then, keep your secrets! Exiting..",
	}
	ErrAPIEndpoint = ErrorType{
		Level:   Error,
		Message: "An error occurred while connecting to the API endpoint",
	}
	ErrPromptFailed = ErrorType{
		Level:   Error,
		Message: "Failed to get context information",
	}
)