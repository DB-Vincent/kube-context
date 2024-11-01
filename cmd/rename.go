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
	"errors"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/pkg/utils"
	"github.com/DB-Vincent/kube-context/pkg/logger"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// Argument definition
var contextFrom string
var contextTo string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Change a context's name",
	Run:   runRenameCommand,
}

// Main logic for rename command
func runRenameCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	// Retrieve contexts and set up configAccess so we can write the adjusted configuration
	opts.GetContexts()
	configAccess := clientcmd.NewDefaultPathOptions()

	// Retrieve context inputs
	validateAndSetContextNames(opts)

	// Rename context
	renameContext(opts, configAccess)
}

func validateAndSetContextNames(opts *utils.KubeConfigOptions) {
	// No contexts were given as argument
	if contextFrom == "" && contextTo == "" {
		promptContextNames(opts)
		return
	}

	if contextFrom == "" || contextTo == "" { // Either "from" or "to" was given, but not both
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Please enter both the name of the context you want to rename and the new name of the context. Use `kube-context rename --help` for more information.",
		}, fmt.Errorf("missing context names"))
		return
	}

	// Verify that context to rename exists in kubeconfig
	if !slices.Contains(opts.Contexts, contextFrom) {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: fmt.Sprintf("Could not find the \"from\" context in kubeconfig file! Found the following contexts: %q", opts.Contexts),
		}, fmt.Errorf("context not found in kubeconfig"))
		return
	}

	// Verify that new name of context doesn't exist in kubeconfig
	_, newExists := opts.Config.Contexts[contextTo]
	if newExists {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "There's already a context with that name. Please give me a different name.",
		}, fmt.Errorf("new context name already exists"))
		return
	}
}

func promptContextNames(opts *utils.KubeConfigOptions) {
	// Set up an interactive prompt to select a context and a new name
	var qs = []*survey.Question{
		{
			Name: "oldContext",
			Prompt: &survey.Select{
				Message: "Choose a context to rename:",
				Options: opts.Contexts,
			},
		},
		{
			Name:     "newContext",
			Prompt:   &survey.Input{Message: "What do you want to name the context?"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		OldContext string `survey:"oldContext"`
		NewContext string `survey:"newContext"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			logHandler.Handle(logger.ErrUserInterrupt, errors.New("user interrupted context rename operation"))
			os.Exit(1)
			return
		} else {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: "Failed to get context information",
			}, err)
			return
		}
	}

	contextFrom = answers.OldContext
	contextTo = answers.NewContext
}

func renameContext(opts *utils.KubeConfigOptions, configAccess clientcmd.ConfigAccess) {
	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Renaming %s context to %s..", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo)),
	}, nil)

	// Get given context
	context, _ := opts.Config.Contexts[contextFrom]

	// Make new context with new name as key and original context as value
	opts.Config.Contexts[contextTo] = context

	// Remove old context
	delete(opts.Config.Contexts, contextFrom)

	// If original context is the current selected context, switch to the new context
	if opts.CurrentContext == contextFrom {
		opts.Config.CurrentContext = contextTo
	}

	// Modify the kubeconfig to ensure that the changes persist
	err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
	if err != nil {
		logHandler.Handle(logger.ErrWriteKubeconfig, err)
		return
	}

	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Successfully renamed %s context to %s!", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo)),
	}, nil)
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&contextFrom, "from", "f", "", "name of context which you want to rename")
	renameCmd.Flags().StringVarP(&contextTo, "to", "t", "", "new name of the context")
}
