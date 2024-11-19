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

var orderCount int64
var re *beep.Resampler
var timeCount int64
var doneAMusic = make(chan bool, 1)
var Loop bool
var loop2 beep.Streamer
var random bool
var order bool
var mL musicLists
var (
	s      beep.StreamSeekCloser
	format beep.Format
)
var count int64
var speakerSampleRate uint32 = 44100
var ctrl *beep.Ctrl
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "play your music",
	Long:  `play your music`,
	Run: func(cmd *cobra.Command, args []string) {
		init := 0
		if random || order {
			sort()
			fmt.Print("Press [ENTER] to pause/resume. ")
			go button()
			for {
				timeCount = rand.Int63()
				open := readMusic(order, random, args)
				s, format, _ = decodeAudioFile(open)
				defer s.Close()
				if init == 0 {
					err := speaker.Init(beep.SampleRate(speakerSampleRate), beep.SampleRate(speakerSampleRate).N(time.Second/10))
					if err != nil {
						fmt.Println(err)
					}
					init = 1
				}
				if init == 1 {
					re = beep.Resample(4, format.SampleRate, beep.SampleRate(speakerSampleRate), s)
				}
				ctrl := switchMode(Loop, random, args, re, init)
				speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
					doneAMusic <- true
				})))
				<-doneAMusic
			}
		} else {
			i := 1
			open := readMusic(order, random, args)
			s, format, _ = decodeAudioFile(open)
			defer s.Close()
			err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			if err != nil {
				fmt.Println(err)
			}
			ctrl := switchMode(Loop, random, args, nil, 0)
			done := make(chan bool, 1)
			speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
				done <- true
			})))
			for i == 1 {
				select {
				case <-done:
					fmt.Println("done")
					i = 0
					return
				default:
					fmt.Print("Press [ENTER] to pause/resume. ")
					fmt.Scanln()
					speaker.Lock()
					ctrl.Paused = !ctrl.Paused
					speaker.Unlock()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.PersistentFlags().BoolVarP(&Loop, "loop", "l", false, "loop the music")
	playCmd.PersistentFlags().BoolVarP(&random, "random", "r", false, "random play your music")
	playCmd.PersistentFlags().BoolVarP(&order, "order", "o", false, "order play your music")
}

func button() {
	for {
		//fmt.Print("Press [ENTER] to pause/resume. ")
		fmt.Scanln()
		speaker.Lock()
		ctrl.Paused = !ctrl.Paused
		speaker.Unlock()
	}
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

func readMusic(order bool, random bool, args []string) *os.File {
	var musicList musicLists
	Db.Model(&musicLists{}).Count(&count)
	if random {
		timeCount++
		rand.New(rand.NewSource(timeCount))
		num := rand.Int63n(count) + 1
		Db.First(&musicList, "id = ?", num)
	} else if order {
		orderCount %= count
		orderCount++
		Db.First(&musicList, "id = ?", orderCount)
	} else {
		var choice string
		if len(args) == 0 {
			choice = "1"
		} else {
			choice = args[0]
		}
		_ = Db.First(&musicList, "id = ?", choice)

	}
	mL = musicList
	fmt.Println(mL.Name)
	open, err := os.Open(mL.Path)
	if err != nil {
		fmt.Println(err)
	}
	return open
}

func switchMode(Loop bool, random bool, args []string, re *beep.Resampler, init int) *beep.Ctrl {
	if Loop {
		var n int
		fmt.Println("loop mode")
		if len(args) <= 1 {
			n = 0
		} else {
			n, _ = strconv.Atoi(args[1])
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
	if init == 1 {
		ctrl = &beep.Ctrl{Streamer: re, Paused: false}
	}
	return ctrl
}
