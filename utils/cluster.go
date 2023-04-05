package utils

import (
	"context"

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
