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
