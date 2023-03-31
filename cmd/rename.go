/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"fmt"

  "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Change a context's name",
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

	  var qs = []*survey.Question{
      {
        Name:   "oldContext",
        Prompt: &survey.Select{
          Message: "Choose a context to rename:",
          Options: contexts,
        },
      },
      {
        Name:     "newContext",
        Prompt:   &survey.Input{Message: "What do you want to name the context?"},
        Validate: survey.Required,
      },
    }

    answers := struct {
      OldContext  string `survey:"oldContext"`
      NewContext  string `survey:"newContext"`
    }{}

    err = survey.Ask(qs, &answers)
    if err != nil {
      if err.Error() == "interrupt" {
        fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
        return
      } else {
        log.Fatal(err.Error())
      }
    }

    context, _ := kubeConfig.Contexts[answers.OldContext]
    _, newExists := kubeConfig.Contexts[answers.NewContext]
    if newExists {
      fmt.Printf("❌ There's already a context with that name. Please give me a different name.\n")
      return
    }

    fmt.Printf("ℹ Renaming %s to %s..\n", answers.OldContext, answers.NewContext)

    kubeConfig.Contexts[answers.NewContext] = context
    delete(kubeConfig.Contexts, answers.OldContext)

    if kubeConfig.CurrentContext == answers.OldContext {
      kubeConfig.CurrentContext = answers.NewContext
    }

    err = clientcmd.ModifyConfig(configAccess, *kubeConfig, true)
    if err != nil {
      log.Fatal("Error %s, modifying config", err.Error())
      return
    }

    fmt.Printf("✔ Successfully renamed %s to %s!\n", answers.OldContext, answers.NewContext)
	},
}

func init() {
  rootCmd.AddCommand(renameCmd)
}
