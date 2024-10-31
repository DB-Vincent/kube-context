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
package utils

import (
	"github.com/DB-Vincent/kube-context/pkg/logger"
)

// Package-level logger variable
var logHandler *logger.Logger

// SetLogger sets the logger instance to be used in this package
func SetLogger(l *logger.Logger) {
	logHandler = l
}