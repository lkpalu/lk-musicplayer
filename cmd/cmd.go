package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lk",
	Short: "lk is a very small music player",
	Long:  `lk is a very small music player,and it's written by golang`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lk is a very small music player")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
