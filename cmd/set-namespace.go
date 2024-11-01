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
var namespace string

// renameCmd represents the rename command
var setDefaultNamespaceCmd = &cobra.Command{
	Use:   "set-namespace",
	Short: "Change a context's default namespace",
	Run:   runSetNamespaceCommand,
}

// Main logic for set-namespace command
func runSetNamespaceCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	// Retrieve contexts and set up configAccess so we can write the adjusted configuration
	opts.GetContexts()
	configAccess := clientcmd.NewDefaultPathOptions()

	// Retrieve namespace to set as default
	selectedNamespace := selectNamespace(opts)
	if selectedNamespace == "" {
		return
	}

	// Sets the namespace
	setNamespace(opts, configAccess, selectedNamespace)
}

func selectNamespace(opts *utils.KubeConfigOptions) string {
	// Retrieve cluster URL, ensuring that we have connection to the cluster
	_, err := opts.GetClusterUrl()
	if err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Failed to connect to the API endpoint",
		}, err)
		return ""
	}

	// Retrieve namespaces in cluster
	if err := opts.GetNamespaces(); err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Failed to retrieve namespaces",
		}, err)
		return ""
	}

	selectedNamespace := ""

	// If no namespace was given, prompt the user to interactively select one
	if namespace == "" {
		selectedNamespace = promptForNamespace(opts)
		if selectedNamespace == "" {
			return ""
		}
	} else { // namespace was given as an argument, verify that it exists in the cluster
		if !slices.Contains(opts.Namespaces, namespace) {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: fmt.Sprintf("Could not find namespace in cluster! Found the following namespaces in your current cluster: %q", opts.Namespaces),
			}, fmt.Errorf("namespace not found in cluster"))
			return ""
		}

		selectedNamespace = namespace
	}

	return selectedNamespace
}

func promptForNamespace(opts *utils.KubeConfigOptions) string {
	result := ""

	// Set up a prompt to interactively select a namespace
	prompt := &survey.Select{
		Message: fmt.Sprintf("Choose a default namespace for the %s context:", opts.CurrentContext),
		Options: opts.Namespaces,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		if err.Error() == "interrupt" {
			logHandler.Handle(logger.ErrUserInterrupt, errors.New("user interrupted namespace selection"))
			os.Exit(0)
			return ""
		} else {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: "Failed to prompt for namespace",
			}, err)
			return ""
		}
	}

	return result
}

func setNamespace(opts *utils.KubeConfigOptions, configAccess clientcmd.ConfigAccess, selectedNamespace string) {
	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Setting the default namespace to %s..", color.FgCyan.Render(selectedNamespace)),
	}, nil)

	// Set namespace parameter for current context
	context, _ := opts.Config.Contexts[opts.CurrentContext]
	context.Namespace = selectedNamespace

	// Write modified configuration to kubeconfig
	if err := clientcmd.ModifyConfig(configAccess, *opts.Config, true); err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Failed to modify kubeconfig",
		}, err)
		return
	}

	logHandler.Handle(logger.ErrorType{
		Level:   logger.Info,
		Message: fmt.Sprintf("Successfully set the default namespace for %s to %s!", color.FgCyan.Render(opts.CurrentContext), color.FgCyan.Render(selectedNamespace)),
	}, nil)
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(setDefaultNamespaceCmd)
	setDefaultNamespaceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "name of namespace you want to set as default")
}
