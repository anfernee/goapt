package cmd

import (
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "apt release related commands",
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
