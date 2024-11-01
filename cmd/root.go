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
	"os"
	"path"
	"slices"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/pkg/utils"
	"github.com/DB-Vincent/kube-context/pkg/logger"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

var (
	debugMode  bool
	logHandler *logger.Logger
)

var rootCmd = &cobra.Command{
	Use:   "kube-context",
	Short: "A simple Go tool to manage Kubernetes contexts in a user-friendly way",
	Long: `kube-context is a command-line interface (CLI) tool designed to simplify the management of Kubernetes contexts, allowing users to seamlessly switch between different Kubernetes clusters with ease.
Whether you are working on multiple projects or interacting with various Kubernetes environments, kube-context provides essential functionality to streamline context management.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize the logger with the debug mode setting
		logHandler = logger.New(debugMode)
		utils.SetLogger(logHandler)
	},
	Run: ContextSwitcher,
}

// Path to kubeconfig file
var kubeConfigPath string

// Context argument
var context string

// Sets the version info for the `kube-context --version` command
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

// Main logic for command
func ContextSwitcher(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	if err := opts.Init(kubeConfigPath); err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Failed to initialize configuration",
		}, err)
		return
	}

	// Retrieve contexts and set up configAccess so we can write the adjusted configuration
	opts.GetContexts()
	configAccess := clientcmd.NewDefaultPathOptions()

	// If no context was given, create an interactive prompt
	if context == "" {
		promptForContext(opts, &context)
		if context == "" {
			return
		}
	} else { // Context argument was given, check if it exists in the kubeconfig file
		if !slices.Contains(opts.Contexts, context) {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: fmt.Sprintf("Could not find context in kubeconfig file! Found the following contexts: %q", opts.Contexts),
			}, fmt.Errorf("context not found"))
			return
		}
	}

	// Switch to the selected context
	switchContext(opts, configAccess, context)
}

func promptForContext(opts *utils.KubeConfigOptions, context *string) {
	// Set up an interactive prompt to select a context
	prompt := &survey.Select{
		Message: "Choose a context:",
		Options: opts.Contexts,
	}

	err := survey.AskOne(prompt, context)
	if err != nil {
		if err.Error() == "interrupt" {
			logHandler.Handle(logger.ErrUserInterrupt, nil)
    	os.Exit(1)
			*context = ""
		} else {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: "Failed to prompt for context",
			}, err)
			*context = ""
		}
	}
}

func switchContext(opts *utils.KubeConfigOptions, configAccess clientcmd.ConfigAccess, context string) {
	// Make sure we're not trying to change to the current context, as that would be pretty pointless
	if opts.CurrentContext != context {
		// Change context to the selected name
		opts.Config.CurrentContext = context

		// Write modified configuration to file
		if err := clientcmd.ModifyConfig(configAccess, *opts.Config, true); err != nil {
			logHandler.Handle(logger.ErrorType{
				Level:   logger.Error,
				Message: "Failed to modify kubeconfig",
			}, err)
			return
		}

		logHandler.Handle(logger.ErrorType{
			Level:   logger.Info,
			Message: fmt.Sprintf("Switched to %s!", color.FgCyan.Render(context)),
		}, nil)
	} else {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Info,
			Message: fmt.Sprintf("You were already working on %s, no need to change.", color.FgCyan.Render(context)),
		}, nil)
	}
}

// Cobra root command caller
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Cobra command initialization
func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		logHandler.Handle(logger.ErrorType{
			Level:   logger.Error,
			Message: "Failed to get user home directory",
		}, err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVarP(&context, "context", "c", "", "name of context to which you want to switch")

	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "config", path.Join(home, ".kube/config"), "kubeconfig file location")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode for detailed logs")
}
