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
	"slices"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
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
	err := validateAndSetContextNames(opts)
	if err != nil {
		fmt.Printf("Error validating and setting context names: %s", err)
		return
	}

	// Rename context
	err = renameContext(opts, configAccess)
	if err != nil {
		fmt.Printf("Error renaming context: %s", err)
		return
	}

	fmt.Printf("✔ Successfully renamed %s context to %s!\n", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo))
}

func validateAndSetContextNames(opts *utils.KubeConfigOptions) error {
	// No contexts were given as argument
	if contextFrom == "" && contextTo == "" {
		err := promptContextNames(opts)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}
	
	if contextFrom == "" || contextTo == "" { // Either "from" or "to" was given, but not both
		fmt.Printf("❌ Please enter both the name of the context you want to rename and the new name of the context. Use `kube-context rename --help` for more information.\n")
		return fmt.Errorf("missing context names")
	}

	// Verify that context to rename exists in kubeconfig
	if !slices.Contains(opts.Contexts, contextFrom) { 
		fmt.Printf("❌ Could not find the \"from\" context in kubeconfig file!\n")
		fmt.Printf("ℹ Found the following contexts in your kubeconfig file: %q\n", opts.Contexts)
		return fmt.Errorf("context not found in kubeconfig")
	}

	// Verify that new name of context doesn't exist in kubeconfig
	_, newExists := opts.Config.Contexts[contextTo]
	if newExists {
		fmt.Printf("❌ There's already a context with that name. Please give me a different name.\n")
		return fmt.Errorf("new context name already exists")
	}

	return nil
}

func promptContextNames(opts *utils.KubeConfigOptions) error {
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
			return fmt.Errorf("ℹ Alright then, keep your secrets! Exiting..\n")
		} else {
			return fmt.Errorf("%s", err)
		}
	}

	contextFrom = answers.OldContext
	contextTo = answers.NewContext

	return nil
}

func renameContext(opts *utils.KubeConfigOptions, configAccess clientcmd.ConfigAccess) error {
	fmt.Printf("ℹ Renaming %s context to %s..\n", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo))

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
	return clientcmd.ModifyConfig(configAccess, *opts.Config, true)
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&contextFrom, "from", "f", "", "name of context which you want to rename")
	renameCmd.Flags().StringVarP(&contextTo, "to", "t", "", "new name of the context")
}
 