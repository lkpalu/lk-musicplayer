package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var mLs []musicLists

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "list your music",
	Long:  `list your music`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = Db.Find(&mLs)
		for _, v := range mLs {
			fmt.Println(v.ID, v.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
