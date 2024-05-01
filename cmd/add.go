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

	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a context to your kubeconfig",
	Run:   runAddCommand,
}

// Main logic for add command
func runAddCommand(cmd *cobra.Command, args []string) {
	var newConfig = false

	// Initialize configuration struct
	opts := &utils.KubeConfigOptions{}
	// Check if kubeconfig file exists
	_, err := os.Stat(kubeConfigPath)
	if os.IsNotExist(err) {
		newConfig = true

		// If kubeconfig file doesn't exist, create an empty config
		opts.Config = &api.Config{
			APIVersion: "v1",
			Kind:       "Config",
			Contexts: map[string]*api.Context{},
			AuthInfos: map[string]*api.AuthInfo{},
			Clusters: map[string]*api.Cluster{},
		}
	} else if err != nil {
		// Error occurred while checking kubeconfig existence
		fmt.Printf("%s", err)
	} else {
		// Load kubeconfig from file
		err := opts.Init(kubeConfigPath)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}
	
	answers := struct {
		Name 		string
		Endpoint 	string
		Authority	string
		Certificate	string
		Key 		string
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

				fmt.Printf("%s", str)

				// if _, exists := opts.Config.Contexts[str]; exists {
				// 	return fmt.Errorf("a context with name '%s' already exists", str)
				// }
				return nil
			},
		},
		{
			Name:     "endpoint",
			Prompt:   &survey.Input{Message: "Please enter the cluster endpoint:"},
			Validate: survey.Required,
		},
		{
			Name:   "authority",
			Prompt: &survey.Input{
				Message: "Please enter the certificate authority location:",
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
	}
	
	err = survey.Ask(prompt, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			fmt.Println("ℹ Alright then, keep your secrets! Exiting..")
		} else {
			fmt.Errorf("%s", err)
		}
	}

	var cluster api.Cluster
	cluster.Server = answers.Endpoint
	cluster.CertificateAuthority = answers.Authority

	var context api.Context
	context.Cluster = answers.Name
	context.AuthInfo = answers.Name

	var auth api.AuthInfo
	auth.ClientCertificate = answers.Certificate
	auth.ClientKey = answers.Key

	opts.Config.Clusters[answers.Name] = &cluster
	opts.Config.Contexts[answers.Name] = &context
	opts.Config.AuthInfos[answers.Name] = &auth

	if (newConfig) {
		err := clientcmd.WriteToFile(*opts.Config, kubeConfigPath)
		if err != nil {
			fmt.Printf("Error saving kubeconfig to file: %v\n", err)
			return
		}
	} else {
		// Write modified configuration to kubeconfig
		configAccess := clientcmd.NewDefaultPathOptions()
		err = clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			fmt.Errorf("error modifying config: %s", err)
			return
		}
	}
}

// Cobra command initialization
func init() {
	rootCmd.AddCommand(addCmd)
}
