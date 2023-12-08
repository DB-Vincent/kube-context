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
	"os"
	"path"
	"slices"

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-context",
	Short: "A simple Go tool to manage Kubernetes contexts in a user-friendly way",
	Long: `kube-context is a command-line interface (CLI) tool designed to simplify the management of Kubernetes contexts, allowing users to seamlessly switch between different Kubernetes clusters with ease. 
Whether you are working on multiple projects or interacting with various Kubernetes environments, kube-context provides essential functionality to streamline context management.`,
	Run: ContextSwitcher,
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

var kubeConfigPath string
var context string

func ContextSwitcher(cmd *cobra.Command, args []string) {
	opts := &utils.KubeConfigOptions{}

	// Initialize environment (retrieve config from file, create clientset)
	opts.Init(kubeConfigPath)
	configAccess := clientcmd.NewDefaultPathOptions()

	// Retrieve contexts from kubeconfig file
	opts.GetContexts()

	result := ""
	if (context == "") {
		prompt := &survey.Select{
			Message: "Choose a context:",
			Options: opts.Contexts,
		}
	
		promptErr := survey.AskOne(prompt, &result)
		if promptErr != nil {
			if promptErr.Error() == "interrupt" {
				fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
				return
			} else {
				log.Fatal(promptErr.Error())
			}
		}	
	} else {
		if (!slices.Contains(opts.Contexts, context)) {
			fmt.Printf("❌ Could not find context in kubeconfig file!\n")
			fmt.Printf("ℹ Found the following contexts in your kubeconfig file: %q\n", opts.Contexts)
			return
		}

		result = context
	}
	
	if opts.CurrentContext != result {
		opts.Config.CurrentContext = result

		err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
		}

		fmt.Printf("✔ Switched to %s!\n", color.FgCyan.Render(result))
	} else {
		fmt.Printf("⚠ You were already working on %s, no need to change.\n", color.FgCyan.Render(result))
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&context, "context", "c", "", "name of context to which you want to switch")

	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "config", path.Join(home, ".kube/config"), "kubeconfig file location")
}
