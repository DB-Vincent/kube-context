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

	"github.com/gookit/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

// deleteCmd represents the rename command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a context from your kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &utils.KubeConfigOptions{}

		// Initialize environment (retrieve config from file, create clientset)
		opts.Init(kubeConfigPath)
		// configAccess := clientcmd.NewDefaultPathOptions()

		// Retrieve contexts from kubeconfig file
		opts.GetContexts()
		configAccess := clientcmd.NewDefaultPathOptions()

		// Selection of context to delete
		var qs = []*survey.Question{
			{
				Name: "contextToDelete",
				Prompt: &survey.Select{
					Message: "Choose a context to delete:",
					Options: opts.Contexts,
				},
			},
		}

		answers := struct {
			ContextToDelete string `survey:"contextToDelete"`
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

		fmt.Printf("ℹ Deleting context %s from kubeconfig file..\n", color.FgCyan.Render(answers.ContextToDelete))

		// Remove context from list of contexts
		delete(opts.Config.Contexts, answers.ContextToDelete)

		var firstContext string

		if opts.CurrentContext == answers.ContextToDelete {
			for context, _ := range opts.Config.Contexts {
        firstContext = context
        break
    }

			fmt.Printf("ℹ You're currently using the context you want to delete, I'll switch you to the %s context..\n", color.FgCyan.Render(firstContext))
			opts.Config.CurrentContext = firstContext
		}

		// Write new context list to kubeconfig file
		err = clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
			return
		}

		fmt.Printf("✔ Successfully deleted context %s!\n", color.FgCyan.Render(answers.ContextToDelete))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
