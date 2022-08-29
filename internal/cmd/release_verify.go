package cmd

import (
	"fmt"
	"os"

	"github.com/anfernee/goapt/pkg/release"
	"github.com/spf13/cobra"
)

var (
	armored    bool
	pubkeyPath string
)

var releaseVerifyCmd = &cobra.Command{
	Use:   "verify <release-path>",
	Short: "Verify apt release file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Missing arguments\n")
			cmd.Usage()
			os.Exit(1)
		}

		path := args[0]
		options := &release.VerifyOptions{}
		if pubkeyPath == "" {
			options.AutoDiscover = true
		} else {
			options.Armored = armored
			options.KeyPath = pubkeyPath
		}

		text, err := release.VerifyWithOptions(path, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(text)
	},
}

func init() {
	releaseCmd.AddCommand(releaseVerifyCmd)

	flags := releaseVerifyCmd.Flags()
	flags.BoolVarP(&armored, "armor", "a", false, "armored public key if specified")
	flags.StringVarP(&pubkeyPath, "pubkey", "k", "", "path to public key")
}
