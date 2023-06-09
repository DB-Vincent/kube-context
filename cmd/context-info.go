/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>

*/
package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/DB-Vincent/kube-context/utils"
	"github.com/spf13/cobra"
)

// contextInfoCmd represents the contextInfo command
var contextInfoCmd = &cobra.Command{
	Use:   "context-info",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &utils.KubeConfigOptions{}

		// Initialize environment (retrieve config from file, create clientset)
		opts.Init(kubeConfigPath)

		// Retrieve and display namespaces
		opts.GetNamespaces()

		// Retrieve cluster url
		clusterUrl, err := opts.GetClusterUrl()
    if err != nil {
      fmt.Printf("❌ An error occurred while connecting to the API endpoint!\nError: %s\n", err.Error())
    }

    fmt.Printf("The %s cluster currently has %s namespaces!\n", opts.Config.CurrentContext, color.FgCyan.Render(len(opts.Namespaces)))
    fmt.Printf("Connecting to this cluster can be done using the %s API endpoint.\n", color.FgCyan.Render(clusterUrl))
	},
}

func init() {
	rootCmd.AddCommand(contextInfoCmd)
}