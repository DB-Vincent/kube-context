/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>
*/
package main

import "kube-context/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
