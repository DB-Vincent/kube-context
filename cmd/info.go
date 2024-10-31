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
	"github.com/DB-Vincent/kube-context/pkg/utils"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve information regarding the current context",
	Run:   runInfoCommand,
}

// Main logic for info command
func runInfoCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	// Retrieves info and displays it to the user
	err := retrieveAndDisplayInfo(opts)
	if err != nil {
		fmt.Printf("Error retrieving and displaying information: %s", err)
		return
	}
}

func retrieveAndDisplayInfo(opts *utils.KubeConfigOptions) error {
	// Retrive cluster URL and make sure that connection works
	clusterUrl, err := opts.GetClusterUrl()
	if err != nil {
		fmt.Printf("❌ An error occurred while connecting to the API endpoint!\nError: %s\n", err.Error())
		return fmt.Errorf("error retrieving cluster URL: %s", err)
	}

	// Retrieve namespaces
	if err := opts.GetNamespaces(); err != nil {
		return fmt.Errorf("error retrieving namespaces: %s", err)
	}

	// Retrieve pods
	if err := opts.GetPods(); err != nil {
		return fmt.Errorf("error retrieving pods: %s", err)
	}

	// Display retrieved information to the user
	displayInfo(opts, clusterUrl)

	return nil
}

func displayInfo(opts *utils.KubeConfigOptions, clusterUrl string) {
	// Print the beautiful data
	fmt.Printf("The %s cluster currently has %s pods spread over %s namespaces!\n", opts.Config.CurrentContext, color.FgCyan.Render(len(opts.Pods)), color.FgCyan.Render(len(opts.Namespaces)))
	fmt.Printf("Connecting to this cluster can be done using the %s API endpoint.\n", color.FgCyan.Render(clusterUrl))
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(infoCmd)
}
