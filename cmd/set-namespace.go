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
	"log"
	"slices"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
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
	selectedNamespace, err := selectNamespace(opts)
	if err != nil {
		log.Fatalf("Error selecting namespace: %s", err)
	}

	// Sets the namespace
	err = setNamespace(opts, configAccess, selectedNamespace)
	if err != nil {
		log.Fatalf("Error setting default namespace: %s", err)
	}
}

func selectNamespace(opts *utils.KubeConfigOptions) (string, error) {
	// Retrieve cluster URL, ensuring that we have connection to the cluster
	_, err := opts.GetClusterUrl()
	if err != nil {
		return "", fmt.Errorf("error connecting to the API endpoint: %s", err)
	}

	// Retrieve namespaces in cluster
	err = opts.GetNamespaces()
	if err != nil {
		return "", fmt.Errorf("error retrieving namespaces: %s", err)
	}

	selectedNamespace := ""

	// If no namespace was given, prompt the user to interactively select one
	if namespace == "" {
		err = promptForNamespace(opts, &selectedNamespace)
		if err != nil {
			return "", fmt.Errorf("error prompting for namespace: %s", err)
		}
	} else { // namespace was given as an argument, verify that it exists in the cluster
		if !slices.Contains(opts.Namespaces, namespace) {
			fmt.Printf("❌ Could not find namespace in cluster!\n")
			fmt.Printf("ℹ Found the following namespaces in your current cluster: %q\n", opts.Namespaces)
			return "", fmt.Errorf("namespace not found in cluster")
		}

		selectedNamespace = namespace
	}

	return selectedNamespace, nil
}

func promptForNamespace(opts *utils.KubeConfigOptions, result *string) error {
	// Set up a prompt to interactively select a namespace
	prompt := &survey.Select{
		Message: fmt.Sprintf("Choose a default namespace for the %s context:", opts.CurrentContext),
		Options: opts.Namespaces,
	}

	err := survey.AskOne(prompt, result)
	if err != nil {
		if err.Error() == "interrupt" {
			fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
			return nil
		} else {
			return fmt.Errorf("error prompting for namespace: %s", err)
		}
	}

	return nil
}


func setNamespace(opts *utils.KubeConfigOptions, configAccess clientcmd.ConfigAccess, selectedNamespace string) error {
	fmt.Printf("ℹ Setting the default namespace to %s..\n", color.FgCyan.Render(selectedNamespace))

	// Set namespace parameter for current context
	context, _ := opts.Config.Contexts[opts.CurrentContext]
	context.Namespace = selectedNamespace

	// Write modified configuration to kubeconfig
	err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
	if err != nil {
		return fmt.Errorf("error modifying config: %s", err)
	}

	fmt.Printf("✔ Successfully set the default namespace for %s to %s!\n", color.FgCyan.Render(opts.CurrentContext), color.FgCyan.Render(selectedNamespace))

	return nil
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(setDefaultNamespaceCmd)
	setDefaultNamespaceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "name of namespace you want to set as default")
}
 