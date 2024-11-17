package cmd

//
//import (
//	"fmt"
//	"github.com/gopxl/beep/v2"
//	"github.com/gopxl/beep/v2/flac"
//	"github.com/gopxl/beep/v2/mp3"
//	"github.com/gopxl/beep/v2/speaker"
//	"github.com/gopxl/beep/v2/wav"
//	"github.com/spf13/cobra"
//
//	//"github.com/spf13/cobra"
//	"math/rand"
//	"os"
//	"path/filepath"
//	"strconv"
//	"time"
//)
//
//var Loop bool
//var random bool
//var mL musicLists
//var ctrl *beep.Ctrl
//var count int64
//var playCmd = &cobra.Command{
//	Use:   "play",
//	Short: "list your music",
//	Long:  `list your music`,
//	Run: func(cmd *cobra.Command, args []string) {
//		playMusic(cmd, args)
//	},
//}
//
//func init() {
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	rootCmd.AddCommand(playCmd)
//	playCmd.PersistentFlags().BoolVarP(&Loop, "loop", "l", false, "loop the music")
//	playCmd.PersistentFlags().BoolVarP(&random, "random", "r", false, "random play your music")
//}
//
//func playMusic(cmd *cobra.Command, args []string) {
//	Db.Model(&musicLists{}).Count(&count)
//	var num int64
//	if random {
//		num = rand.Int63n(count) + 1
//	} else {
//		n, _ := strconv.Atoi(args[0])
//		num = int64(n)
//	}
//	_ = Db.First(&mL, "id = ?", num)
//	fmt.Println(mL.Name)
//
//	open, err := os.Open(mL.Path)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer open.Close()
//
//	s, format, err := decodeAudioFile(open)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	ctrl = &beep.Ctrl{Streamer: s, Paused: false}
//	if Loop {
//		n := 1
//		if len(args) > 1 {
//			n, err = strconv.Atoi(args[1])
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//		}
//		if n == 0 {
//			ctrl.Streamer, _ = beep.Loop2(s)
//		} else {
//			ctrl.Streamer, _ = beep.Loop2(s, beep.LoopTimes(n-1))
//		}
//	}
//
//	done := make(chan bool, 1)
//	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
//		done <- true
//	})))
//
//	playControl(done)
//}
//
//func decodeAudioFile(file *os.File) (beep.StreamSeekCloser, beep.Format, error) {
//	ext := filepath.Ext(file.Name())
//	switch ext {
//	case ".mp3":
//		return mp3.Decode(file)
//	case ".wav":
//		return wav.Decode(file)
//	case ".flac":
//		return flac.Decode(file)
//	default:
//		return nil, beep.Format{}, fmt.Errorf("unsupported file format: %s", ext)
//	}
//}
//
//func playControl(done chan bool) {
//	for {
//		fmt.Print("Press [ENTER] to pause/resume. ")
//		fmt.Scanln()
//		speaker.Lock()
//		ctrl.Paused = !ctrl.Paused
//		speaker.Unlock()
//		if random && <-done {
//			fmt.Println("done")
//			return
//		}
//	}
//}
