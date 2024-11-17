package cmd

import (
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "sort your music list",
	Long:  `sort your music list`,
	Run: func(cmd *cobra.Command, args []string) {
		var users []musicLists
		Db.Find(&users)

		for i, user := range users {
			Db.Model(&user).Where("id = ?", user.ID).Update("id", i+1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}
