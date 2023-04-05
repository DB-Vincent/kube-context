/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-context",
	Short: "A simple Go tool to switch between Kubernetes contexts in a user-friendly way",
	Long: `Kube-context is a simple and easy-to-use CLI tool written in Go,
which allows you to choose a Kubernetes config in a user-friendly way.

It simplifies the process of switching between Kubernetes contexts by providing
a menu-driven interface to list, select and switch between Kubernetes contexts.`,
	Run: ContextSwitcher,
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

var kubeConfigPath string

func ContextSwitcher(cmd *cobra.Command, args []string) {
	opts := &utils.KubeConfigOptions{}

	// Initialize environment (retrieve config from file, create clientset)
	opts.Init(kubeConfigPath)
	configAccess := clientcmd.NewDefaultPathOptions()

	// Retrieve contexts from kubeconfig file
	opts.GetContexts()

	result := ""
	prompt := &survey.Select{
		Message: "Choose a context:",
		Options: opts.Contexts,
	}

	promptErr := survey.AskOne(prompt, &result)
	if promptErr != nil {
		if promptErr.Error() == "interrupt" {
			fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
			return
		} else {
			log.Fatal(promptErr.Error())
		}
	}

	if opts.CurrentContext != result {
		opts.Config.CurrentContext = result

		err := clientcmd.ModifyConfig(configAccess, *opts.Config, true)
		if err != nil {
			log.Fatal("Error %s, modifying config", err.Error())
		}

		fmt.Printf("✔ Switched to %s!\n", result)
	} else {
		fmt.Printf("⚠ You were already working on %s, no need to change.\n", result)
	}
}

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
