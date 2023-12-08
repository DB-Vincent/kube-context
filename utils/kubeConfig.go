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

	opts.Config, err = clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	opts.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func (opts *KubeConfigOptions) GetContexts() {
	for context := range opts.Config.Contexts {
		opts.Contexts = append(opts.Contexts, context)
	}

	opts.CurrentContext = opts.Config.CurrentContext
}
