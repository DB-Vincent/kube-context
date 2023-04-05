/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

// setDefaultNamespaceCmd represents the setNamespace command
var setDefaultNamespaceCmd = &cobra.Command{
	Use:   "set-default-namespace",
	Short: "Change a context's default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &utils.KubeConfigOptions{}

		// Initialize environment (retrieve config from file, create clientset)
		opts.Init(kubeConfigPath)
		configAccess := clientcmd.NewDefaultPathOptions()

		// Retrieve contexts from kubeconfig file
		opts.GetContexts()

		// Retrieve namespaces for current context
		err := opts.GetNamespaces()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Retrieve connection URL and test connectivity
		currentClusterName := opts.Config.Contexts[opts.CurrentContext].Cluster
		connectionUrl := opts.Config.Clusters[currentClusterName].Server

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

		// Display namespace selection prompt to user
		selectedNamespace := ""
		prompt := &survey.Select{
			Message: fmt.Sprintf("Choose a default namespace for the \"%s\" context:", opts.CurrentContext),
			Options: opts.Namespaces,
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

		// Change namespace in kubeconfig to selected namespace
		fmt.Printf("ℹ Setting the default namespace to \"%s\"..\n", selectedNamespace)
		context, _ := opts.Config.Contexts[opts.CurrentContext]
		context.Namespace = selectedNamespace
		err = clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
			return
		}

		fmt.Printf("✔ Successfully set the default namespace for \"%s\" to \"%s\"!\n", opts.CurrentContext, selectedNamespace)
	},
}

func init() {
	rootCmd.AddCommand(setDefaultNamespaceCmd)
}
