package cmd

import (
	"fmt"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var Loop bool
var loop2 beep.Streamer
var random bool
var mL musicLists
var (
	s      beep.StreamSeekCloser
	format beep.Format
)
var count int64

var ctrl *beep.Ctrl
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "play your music",
	Long:  `play your music`,
	Run: func(cmd *cobra.Command, args []string) {
		if random {
			Db.Model(&musicLists{}).Count(&count)
			num := rand.Int63n(count-1+1) + 1
			_ = Db.First(&mL, "id = ?", num)
		} else {
			_ = Db.First(&mL, "id = ?", args[0])
		}
		fmt.Println(mL.Name)
		open, err := os.Open(mL.Path)
		if err != nil {
			fmt.Println(err)
		}
		s, format, _ = decodeAudioFile(open)
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
		done := make(chan bool, 1)
		speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
			done <- true
		})))
		for {
			select {
			case <-done:
				fmt.Println("done")
				return
			default:
				//fmt.Print(done)
				fmt.Print("Press [ENTER] to pause/resume. ")
				fmt.Scanln()
				speaker.Lock()
				ctrl.Paused = !ctrl.Paused
				speaker.Unlock()
			}
		}
	},
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rootCmd.AddCommand(playCmd)
	playCmd.PersistentFlags().BoolVarP(&Loop, "loop", "l", false, "loop the music")
	playCmd.PersistentFlags().BoolVarP(&random, "random", "r", false, "random play your music")
}

func decodeAudioFile(file *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	ext := filepath.Ext(file.Name())
	switch ext {
	case ".mp3":
		return mp3.Decode(file)
	case ".wav":
		return wav.Decode(file)
	case ".flac":
		return flac.Decode(file)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported file format: %s", ext)
	}
}
