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

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api"
	"github.com/DB-Vincent/kube-context/pkg/logger"
)

type KubeConfigOptions struct {
	Namespaces     []string
	Pods					 []string
	Contexts       []string
	CurrentContext string

	Config *api.Config
	Client *kubernetes.Clientset
}

func (opts *KubeConfigOptions) Init(kubeConfigPath string) {
	// Load kube config file
	var err error
	opts.Config, err = clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		logHandler.Handle(logger.ErrInitKubeconfig, err)
		return
	}

	// Build client-usable configuration from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logHandler.Handle(logger.ErrAPIEndpoint, err)
		return
	}

	// Create client from previously retrieved configuration
	opts.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		logHandler.Handle(logger.ErrAPIEndpoint, err)
		return
	}
}

func (opts *KubeConfigOptions) InitOrCreate(kubeConfigPath string) {
	_, err := os.Stat(kubeConfigPath)
	if os.IsNotExist(err) {
		// If kubeconfig file doesn't exist, create an empty config
		opts.Config = &api.Config{
			APIVersion: "v1",
			Kind:       "Config",
			Contexts:   map[string]*api.Context{},
			AuthInfos:  map[string]*api.AuthInfo{},
			Clusters:   map[string]*api.Cluster{},
		}

		// No Kubeconfig was present, so we create one with the new data
		err := clientcmd.WriteToFile(*opts.Config, kubeConfigPath)
		if err != nil {
			logHandler.Handle(logger.ErrWriteKubeconfig, err)
			return
		}
	} else if err != nil {
		logHandler.Handle(logger.ErrInitKubeconfig, err)
		return
	} else {
		// Load kubeconfig from file
		opts.Init(kubeConfigPath)
	}
}

func (opts *KubeConfigOptions) GetContexts() {
	// Loop through contexts inside kubeconfig file
	for context := range opts.Config.Contexts {
		opts.Contexts = append(opts.Contexts, context)
	}

	opts.CurrentContext = opts.Config.CurrentContext
}
