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

var contextFrom string
var contextTo string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Change a context's name",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &utils.KubeConfigOptions{}

		// Initialize environment (retrieve config from file, create clientset)
		opts.Init(kubeConfigPath)
		configAccess := clientcmd.NewDefaultPathOptions()

		// Retrieve contexts from kubeconfig file
		opts.GetContexts()

		if (contextFrom == "" && contextTo == "") {
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
					fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
					return
				} else {
					log.Fatal(err.Error())
				}
			}
	
			contextFrom = answers.OldContext
			contextTo = answers.NewContext
		} else if (contextFrom == "" || contextTo == "") {
			fmt.Printf("❌ Please give enter both the name of the context you want to rename and the new name of the context. Use `kube-context rename --help` for more information.\n")
			return
		} else {
			if (!slices.Contains(opts.Contexts, contextFrom)) {
				fmt.Printf("❌ Could not find the \"from\" context in kubeconfig file!\n")
				fmt.Printf("ℹ Found the following contexts in your kubeconfig file: %q\n", opts.Contexts)
				return
			}
		}

		context, _ := opts.Config.Contexts[contextFrom]
		_, newExists := opts.Config.Contexts[contextTo]
		if newExists {
			fmt.Printf("❌ There's already a context with that name. Please give me a different name.\n")
			return
		}

		fmt.Printf("ℹ Renaming %s context to %s..\n", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo))

		opts.Config.Contexts[contextTo] = context
		delete(opts.Config.Contexts, contextFrom)

		if opts.CurrentContext == contextFrom {
			opts.Config.CurrentContext = contextTo
		}

		err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
			return
		}

		fmt.Printf("✔ Successfully renamed %s context to %s!\n", color.FgCyan.Render(contextFrom), color.FgCyan.Render(contextTo))
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	renameCmd.Flags().StringVarP(&contextFrom, "from", "f", "", "Name of context which you want to rename")
	renameCmd.Flags().StringVarP(&contextTo, "to", "t", "", "New name of the context")
}
