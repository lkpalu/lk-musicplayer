package cmd

import (
	"fmt"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var mL musicLists
var (
	s      beep.StreamSeekCloser
	format beep.Format
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "play your music",
	Long:  `play your music`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = Db.First(&mL, "id = ?", args[0])
		fmt.Println(mL.Name)
		open, err := os.Open(mL.Path)
		if err != nil {
			fmt.Println(err)
		}
		ext := filepath.Ext(mL.Name)
		switch ext {
		case ".mp3":
			s, format, _ = mp3.Decode(open)
		case ".wav":
			s, format, _ = wav.Decode(open)
		case ".flac":
			s, format, _ = flac.Decode(open)

		}

		if err != nil {
			fmt.Println(err)
		}
		defer s.Close()
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			fmt.Println(err)
		}
		done := make(chan bool)
		speaker.Play(beep.Seq(s, beep.Callback(func() {
			done <- true
		})))

		<-done
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}
