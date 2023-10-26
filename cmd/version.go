package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本号",
	Long:  `打印版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s \n", "0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
