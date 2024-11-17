package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var All bool
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove your music",
	Long:  `remove your music`,
	Run: func(cmd *cobra.Command, args []string) {
		if All {
			err := Db.Migrator().DropTable(&musicLists{})
			if err != nil {
				fmt.Println("delete failed", err)
			}
			return
		}
		num, _ := strconv.Atoi(args[0])
		Db.Unscoped().Delete(&musicLists{}, num)
	},
}

func init() {
	rmCmd.PersistentFlags().BoolVarP(&All, "all", "a", false, "remove all music")
	rootCmd.AddCommand(rmCmd)
}
