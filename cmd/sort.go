package cmd

import (
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "sort your music list",
	Long:  `sort your music list`,
	Run: func(cmd *cobra.Command, args []string) {
		var musicList []musicLists
		Db.Find(&musicList)

		for i, music := range musicList {
			Db.Model(&music).Where("id = ?", music.ID).Update("id", i+1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}
