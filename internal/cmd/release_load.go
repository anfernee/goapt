package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var releaseLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load apt release file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: load")
	},
}

func init() {
	releaseCmd.AddCommand(releaseLoadCmd)
}
