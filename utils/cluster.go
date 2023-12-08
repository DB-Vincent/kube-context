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
	"errors"
	"crypto/tls"
  "net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (opts *KubeConfigOptions) GetNamespaces() error {
	var err error

	namespaceList, err := opts.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, n := range namespaceList.Items {
		opts.Namespaces = append(opts.Namespaces, n.Name)
	}

	return nil
}

func (opts *KubeConfigOptions) GetPods() error {
	var err error

	podList, err := opts.Client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range podList.Items {
		opts.Pods = append(opts.Pods, pod.Name)
	}

	return nil
}

func (opts *KubeConfigOptions) GetClusterUrl() (string, error) {
	var err error

	// Retrieve connection URL and test connectivity
  currentClusterName := opts.Config.Contexts[opts.Config.CurrentContext].Cluster
  connectionUrl := opts.Config.Clusters[currentClusterName].Server

  http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

  response, err := http.Get(connectionUrl)
  if err != nil {
    return "", err
  } else {
    if response.StatusCode != 401 { // We can expect to be hit with an "Unauthorized" message, this *should* be fine.
      return "", errors.New("Did not receive expected \"401\" HTTP status code!")
    }
  }

  response.Body.Close()

  return connectionUrl, nil
}