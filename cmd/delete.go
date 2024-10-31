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
	"os"
	"fmt"
	"slices"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/pkg/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a context from your kubeconfig",
	Run:   runDeleteCommand,
}

// Main logic for delete command
func runDeleteCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	// Retrieve contexts and set up configAccess so we can write the adjusted configuration
	opts.GetContexts()
	configAccess := clientcmd.NewDefaultPathOptions()

	// Prompt the user to select a context to delete
	contextToDelete, err := selectContextToDelete(opts)
	if err != nil {
		fmt.Printf("Error selecting context to delete: %s", err)
		return
	}

	// Remove selected context from kubeconfig
	err = deleteContext(opts, contextToDelete)
	if err != nil {
		fmt.Printf("Error deleting context: %s", err)
		return
	}

	// Write modified configuration to kubeconfig file
	err = clientcmd.ModifyConfig(configAccess, *opts.Config, true)
	if err != nil {
		fmt.Printf("Error updating kubeconfig file: %s", err)
		return
	}

	fmt.Printf("✔ Successfully deleted context %s!\n", color.FgCyan.Render(contextToDelete))
}

func selectContextToDelete(opts *utils.KubeConfigOptions) (string, error) {
	// If a context was given as an argument, check if it exists in the kubeconfig
	if context != "" {
		if !slices.Contains(opts.Contexts, context) {
			fmt.Printf("❌ Could not find context in kubeconfig file!\n")
			fmt.Printf("ℹ Found the following contexts in your kubeconfig file: %q\n", opts.Contexts)
			return "", fmt.Errorf("context not found in kubeconfig")
		}
		return context, nil
	}

	// No context was given, set up a prompt to interactively select context
	prompt := &survey.Select{
		Message: "Choose a context to delete:",
		Options: opts.Contexts,
	}

	err := survey.AskOne(prompt, &context)
	if err != nil {
		if err.Error() == "interrupt" {
			fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
			os.Exit(1)
			return "", nil
		} else {
			return "", fmt.Errorf("error selecting context: %s", err)
		}
	}

	return context, nil
}

func deleteContext(opts *utils.KubeConfigOptions, contextToDelete string) error {
	fmt.Printf("ℹ Deleting context %s from kubeconfig file..\n", color.FgCyan.Render(contextToDelete))

	// Remove context from context list in configuration struct
	delete(opts.Config.Contexts, contextToDelete)

	// Change current context to first context in list if current context is deleted
	if opts.CurrentContext == contextToDelete {
		firstContext := getFirstContext(opts.Config.Contexts)
		fmt.Printf("ℹ You're currently using the context you want to delete, I'll switch you to the %s context..\n", color.FgCyan.Render(firstContext))
		opts.Config.CurrentContext = firstContext
	}

	return nil
}

func getFirstContext(contexts map[string]*api.Context) string {
	var firstContext string

	// Loop through contexts and return first item
	for context := range contexts {
		firstContext = context
		break
	}
	return firstContext
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&context, "context", "c", "", "name of context which you want to delete")
}
