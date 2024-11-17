package cmd

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type musicLists struct {
	gorm.Model
	Name string
	Path string
}

var Db *gorm.DB

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add your music",
	Long:  `add your music`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		readDir, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
		}
		//i := 0
		for _, file := range readDir {
			ext := filepath.Ext(file.Name())
			if ext == ".mp3" || ext == ".flac" || ext == ".wav" {
				join := filepath.Join(dir, file.Name())
				//fmt.Println(i, " ", file.Name())
				//db.Model(&musicLists{}).Update(s, file.Name())
				Db.Create(&musicLists{Name: file.Name(), Path: join})
				//i++
			}
		}
	},
}

type Config struct {
	Root string `json:"root"`
}

func init() {
	configPath := os.Getenv("Mytool")
	//fmt.Println(configPath)
	file, err := os.Open(fmt.Sprintf("%s/config.json", configPath))
	fmt.Println(fmt.Sprintf("%s/config.json", configPath))
	if err != nil {
		panic("打开配置文件失败")
	}
	defer file.Close()
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	Db, _ = gorm.Open(sqlite.Open(fmt.Sprintf("%s/musicPlayer.db", config.Root)), &gorm.Config{})
	_ = Db.AutoMigrate(&musicLists{})
	rootCmd.AddCommand(addCmd)
}
