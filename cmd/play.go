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
	"strconv"
	"time"
)

var Loop bool
var loop2 beep.Streamer
var mL musicLists
var (
	s      beep.StreamSeekCloser
	format beep.Format
)
var ctrl *beep.Ctrl
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
		if Loop {
			fmt.Println("loop mode")
			n, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
			}
			if n == 0 {
				loop2, _ = beep.Loop2(s)
			} else {
				loop2, _ = beep.Loop2(s, beep.LoopTimes(n-1))
			}
			ctrl = &beep.Ctrl{Streamer: loop2, Paused: false}
		} else {
			ctrl = &beep.Ctrl{Streamer: s, Paused: false}
		}
		done := make(chan bool)
		//speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		//	done <- true
		//})))
		speaker.Play(ctrl)
		for {
			fmt.Print("Press [ENTER] to pause/resume. ")
			fmt.Scanln()
			speaker.Lock()
			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
			err := ctrl.Err()
			if err != nil {
				fmt.Println(err)
				done <- true
				break
			}

		}
		<-done
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.PersistentFlags().BoolVarP(&Loop, "loop", "l", false, "loop the music")
}
