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
	"github.com/DB-Vincent/kube-context/pkg/logger"
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
	contextToDelete := selectContextToDelete(opts)
	if contextToDelete == "" {
		return
	}

	// Remove selected context from kubeconfig
	deleteContext(opts, contextToDelete)

	// Write modified configuration to kubeconfig file
	err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
	if err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Error updating kubeconfig file",
		}, err)
		return
	}

	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Successfully deleted context %s!", color.FgCyan.Render(contextToDelete)),
	}, nil)
}

func selectContextToDelete(opts *utils.KubeConfigOptions) string {
	// If a context was given as an argument, check if it exists in the kubeconfig
	if context != "" {
		if !slices.Contains(opts.Contexts, context) {
			logHandler.Handle(logger.ErrContextNotFound, fmt.Errorf("context not found in kubeconfig"), opts.Contexts)
			return ""
		}
		return context
	}

	// No context was given, set up a prompt to interactively select context
	prompt := &survey.Select{
		Message: "Choose a context to delete:",
		Options: opts.Contexts,
	}

	err := survey.AskOne(prompt, &context)
	if err != nil {
		if err.Error() == "interrupt" {
			logHandler.Handle(logger.ErrUserInterrupt, nil)
			os.Exit(1)
			return ""
		} else {
			logHandler.Handle(logger.ErrSelectContext, err)
			return ""
		}
	}

	return context
}

func deleteContext(opts *utils.KubeConfigOptions, contextToDelete string) {
	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Deleting context %s from kubeconfig file..", color.FgCyan.Render(contextToDelete)),
	}, nil)

	// Remove context from context list in configuration struct
	delete(opts.Config.Contexts, contextToDelete)

	// Change current context to first context in list if current context is deleted
	if opts.CurrentContext == contextToDelete {
		firstContext := getFirstContext(opts.Config.Contexts)
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Info,
			Message: fmt.Sprintf("You're currently using the context you want to delete, I'll switch you to the %s context..", color.FgCyan.Render(firstContext)),
		}, nil)
		opts.Config.CurrentContext = firstContext
	}
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
