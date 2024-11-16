package cmd

import (
	"github.com/spf13/cobra"
	"strconv"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove your music",
	Long:  `remove your music`,
	Run: func(cmd *cobra.Command, args []string) {
		num, _ := strconv.Atoi(args[0])
		Db.Delete(&musicLists{}, num)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
