/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"os"
	"path"
	"log"
  "fmt"

  "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfigPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-context",
	Short: "A simple Go tool to switch between Kubernetes contexts in a user-friendly way",
	Long: `Kube-context is a simple and easy-to-use CLI tool written in Go,
which allows you to choose a Kubernetes config in a user-friendly way.

It simplifies the process of switching between Kubernetes contexts by providing
a menu-driven interface to list, select and switch between Kubernetes contexts.`,

	Run: func(cmd *cobra.Command, args []string) {
    kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
    configAccess := clientcmd.NewDefaultPathOptions()
    contexts := []string{}
    if err != nil {
      log.Fatal(err)
    }

    for name := range kubeConfig.Contexts {
      contexts = append(contexts, name,)
    }

    result := ""
    prompt := &survey.Select{
        Message: "Choose a context:",
        Options: contexts,
    }
    survey.AskOne(prompt, &result)

    if kubeConfig.CurrentContext != result {
      kubeConfig.CurrentContext = result

      err = clientcmd.ModifyConfig(configAccess, *kubeConfig, true)
      if err != nil {
        log.Fatal("Error %s, modifying config", err.Error())
      }

      fmt.Printf("✔ Switched to %s!\n", result)
    } else {
      fmt.Printf("⚠ You were already working on %s, no need to change.\n", result)
    }
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
  home, err := os.UserHomeDir()
  if err != nil {
    log.Fatal(err)
  }

	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "config", path.Join(home, ".kube/config"), "Kubeconfig file location")
}


