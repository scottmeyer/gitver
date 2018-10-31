package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configuration utility for gitversion",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called.")
	},
}
