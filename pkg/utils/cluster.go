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
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/DB-Vincent/kube-context/pkg/logger"
)

// GetNamespaces retrieves a list of namespaces in the current cluster.
func (opts *KubeConfigOptions) GetNamespaces() {
	namespaceList, err := opts.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logHandler.Handle(logger.ErrGetResource, err, "namespace")
	}

	for _, n := range namespaceList.Items {
		opts.Namespaces = append(opts.Namespaces, n.Name)
	}
}

// GetPods retrieves a list of pods in the current cluster.
func (opts *KubeConfigOptions) GetPods() {
	podList, err := opts.Client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logHandler.Handle(logger.ErrGetResource, err, "pod")
	}

	for _, pod := range podList.Items {
		opts.Pods = append(opts.Pods, pod.Name)
	}
}

// GetClusterUrl retrieves the connection URL of the current cluster and tests connectivity.
func (opts *KubeConfigOptions) GetClusterUrl() string {
	currentClusterName := opts.Config.Contexts[opts.Config.CurrentContext].Cluster
	connectionURL := opts.Config.Clusters[currentClusterName].Server

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	response, err := http.Get(connectionURL)
	if err != nil {
		logHandler.Handle(logger.ErrAPIEndpoint, err)
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		logHandler.Handle(logger.ErrAPIEndpoint, errors.New("did not receive expected \"401\" HTTP status code"))
		return ""
	}

	return connectionURL
}
