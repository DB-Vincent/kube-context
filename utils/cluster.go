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