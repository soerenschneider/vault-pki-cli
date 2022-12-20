package main

import (
	"fmt"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.BuildVersion)
	},
}
