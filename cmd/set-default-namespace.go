/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// setDefaultNamespaceCmd represents the setNamespace command
var setDefaultNamespaceCmd = &cobra.Command{
	Use:   "set-default-namespace",
	Short: "Change a context's default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
      log.Fatal(err.Error())
    }

		kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
		if err != nil {
      log.Fatal(err.Error())
    }

		configAccess := clientcmd.NewDefaultPathOptions()
		namespaces := []string{}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(err.Error())
		}

		currentClusterName := kubeConfig.Contexts[kubeConfig.CurrentContext].Cluster
		connectionUrl := kubeConfig.Clusters[currentClusterName].Server

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		response, err := http.Get(connectionUrl)
		if err != nil {
			fmt.Printf("❌ An error occurred while connecting to the API endpoint for \"%s\" (%s)!\nError: %s\n", currentClusterName, connectionUrl, err.Error())
			return
		} else {
			if response.StatusCode != 401 { // We can expect to be hit with an "Unauthorized" message, this *should* be fine.
				fmt.Printf("❌ Couldn't connect to the API endpoint for \"%s\" (%s)!\n", currentClusterName, connectionUrl)
				return
			}
		}
		response.Body.Close()

		namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, n := range namespaceList.Items {
			namespaces = append(namespaces, n.Name)
		}

		selectedNamespace := ""
		prompt := &survey.Select{
			Message: fmt.Sprintf("Choose a default namespace for the \"%s\" context:", kubeConfig.CurrentContext),
			Options: namespaces,
		}

		promptErr := survey.AskOne(prompt, &selectedNamespace)
		if promptErr != nil {
			if promptErr.Error() == "interrupt" {
				fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
				return
			} else {
				log.Fatal(promptErr.Error())
			}
		}

		fmt.Printf("ℹ Setting the default namespace to \"%s\"..\n", selectedNamespace)
		context, _ := kubeConfig.Contexts[kubeConfig.CurrentContext]
		context.Namespace = selectedNamespace
		err = clientcmd.ModifyConfig(configAccess, *kubeConfig, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
			return
		}

		fmt.Printf("✔ Successfully set the default namespace for \"%s\" to \"%s\"!\n", kubeConfig.CurrentContext, selectedNamespace)
	},
}

func init() {
	rootCmd.AddCommand(setDefaultNamespaceCmd)
}
