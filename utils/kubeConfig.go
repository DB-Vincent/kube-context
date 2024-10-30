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
package utils

import (
	"os"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfigOptions struct {
	Namespaces     []string
	Pods					 []string
	Contexts       []string
	CurrentContext string

	Config *api.Config
	Client *kubernetes.Clientset
}

func (opts *KubeConfigOptions) Init(kubeConfigPath string) error {
	var err error

	// Load kube config file
	opts.Config, err = clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return err
	}

	// Build client-usable configuration from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	// Create client from previously retrieved configuration
	opts.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func (opts *KubeConfigOptions) InitOrCreate(kubeConfigPath string) error {
	var err error

	_, err = os.Stat(kubeConfigPath)
	if os.IsNotExist(err) {
		// If kubeconfig file doesn't exist, create an empty config
		opts.Config = &api.Config{
			APIVersion: "v1",
			Kind:       "Config",
			Contexts: map[string]*api.Context{},
			AuthInfos: map[string]*api.AuthInfo{},
			Clusters: map[string]*api.Cluster{},
		}

		// No Kubeconfig was present, so we create one with the new data
		err := clientcmd.WriteToFile(*opts.Config, kubeConfigPath)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Error occurred while checking kubeconfig existence
		fmt.Printf("%v\n", err)
		return err
	} else {
		// Load kubeconfig from file
		err := opts.Init(kubeConfigPath)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

	return nil
}

func (opts *KubeConfigOptions) GetContexts() {
	// Loop through contexts inside kubeconfig file
	for context := range opts.Config.Contexts {
		opts.Contexts = append(opts.Contexts, context)
	}

	opts.CurrentContext = opts.Config.CurrentContext
}
