/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"os"
	"path"
	"log"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/manifoldco/promptui"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-context",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
	  home, err := os.UserHomeDir()
    kubeConfig, err := clientcmd.LoadFromFile(path.Join(home, ".kube/config"))
    configAccess := clientcmd.NewDefaultPathOptions()
    contexts := []string{}
    if err != nil {
      log.Fatal(err)
    }

    for name := range kubeConfig.Contexts {
      contexts = append(contexts, name,)
    }

    prompt := promptui.Select{
      Label: "Select context",
      Items: contexts,
    }

    _, result, err := prompt.Run()

    if err != nil {
      log.Fatal("Prompt failed %v\n", err)
      return
    }

    kubeConfig.CurrentContext = result
    err = clientcmd.ModifyConfig(configAccess, *kubeConfig, true)
    if err != nil {
      log.Fatal("Error %s, modifying config", err.Error())
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-context.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
// 	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


