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
	"errors"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a context to your kubeconfig",
	Run:   runAddCommand,
}

// Main logic for add command
func runAddCommand(cmd *cobra.Command, args []string) {
	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	opts.Init(kubeConfigPath)

	answers := struct {
		Name 		string
		Endpoint 	string
		Certificate	string
		Key 		string
		Continue	bool
	}{}

	// Set up an interactive prompt to select a context
	var prompt = []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Please enter a name for the context:"},
			Validate: func (val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("input value is not a string")
				}

				if _, exists := opts.Config.Contexts[str]; exists {
					return fmt.Errorf("a context with name '%s' already exists", str)
				}
				return nil
			},
		},
		{
			Name:     "endpoint",
			Prompt:   &survey.Input{Message: "Please enter the cluster endpoint:"},
			Validate: survey.Required,
		},
		{
			Name:   "certificate",
			Prompt: &survey.Input{
				Message: "Please enter the client certificate location:",
				Suggest: func (toComplete string) []string {
					files, _ := filepath.Glob(toComplete + "*")
					return files
				},
			},
			Validate: func (val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("input value is not a string")
				}

				if _, err := os.Stat(str); errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("could not find a file with name '%s'", str)
				}
				return nil
			},
		},
		{
			Name:   "key",
			Prompt: &survey.Input{
				Message: "Please enter the client key location:",
				Suggest: func (toComplete string) []string {
					files, _ := filepath.Glob(toComplete + "*")
					return files
				},
			},
			Validate: func (val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("input value is not a string")
				}

				if _, err := os.Stat(str); errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("could not find a file with name '%s'", str)
				}
				return nil
			},
		},
		{
			Name:  "continue",
			Prompt: &survey.Confirm{
				Message: "Does the information above look right?",
			},
		},
	}
	
	err := survey.Ask(prompt, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			fmt.Println("ℹ Alright then, keep your secrets! Exiting..")
		} else {
			fmt.Errorf("%s", err)
		}
	}

	if (!answers.Continue) {
		fmt.Println("ℹ Hm, okay. Can you try launching the command again and providing the right information?\n")
		return
	}

	fmt.Printf("name: %s\nendpoint: %s\ncertificate: %s\nkey: %s", answers.Name, answers.Endpoint, answers.Certificate, answers.Key)
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(addCmd)
}
