/*
 * kube-context
 *
 * Copyright (C) 2023 Vincent De Borger
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
package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available contexts in kubeconfig",
	Run:   runListCommand,
}

// Main logic for list command
func runListCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	// Retrieve contexts
	opts.GetContexts()

	fmt.Printf("You currently have %s context(s) configured:\n", color.FgCyan.Render(len(opts.Contexts)))
	for _, context := range opts.Contexts {
		fmt.Printf("- %s\n", color.FgCyan.Render(context))
	}
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(listCmd)
}
