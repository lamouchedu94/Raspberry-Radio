package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

type info struct {
	stationName string
	stationFreq int
	stationText any
}

func main() {
	i := info{}
	// Create a context with cancel function to gracefully handle Ctrl+C events

	ctx, cancel := context.WithCancel(context.Background())
	// Handle Ctrl+C signal (SIGINT)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		fmt.Println("\nCtrl+C received. shutting down...")
		//cmd_kill := exec.Command("bash", "-c", "kill $(ps -e -o pid,command | grep 'InputAndScreen.py' | awk '{print $1}')")
		//cmd_kill.Run()

		//kill $(ps -e -o pid,comm | grep 'InputAndScreen.py' | awk '{print $1}')

		cancel() // Cancel the context when Ctrl+C is received
		os.Exit(1)
	}()
	err := i.run(ctx)
	if err != nil {
		fmt.Println(err)
	}

}

func (i *info) pythonInterface(cmd_python *exec.Cmd) error {

	pipeIn, _ := cmd_python.StdinPipe()
	pipeOut, _ := cmd_python.StdoutPipe()

	err := cmd_python.Start()
	if err != nil {
		return err

	}

	go func() {
		old_name := i.stationName
		old_freq := i.stationFreq
		data := fmt.Sprintf("freq %d\n", i.stationFreq)
		pipeIn.Write([]byte(data))
		data = fmt.Sprintf("name %s\n", i.stationName)
		pipeIn.Write([]byte(data))
		for {
			if old_freq != i.stationFreq {
				data := fmt.Sprintf("freq %d\n", i.stationFreq)
				pipeIn.Write([]byte(data))
				old_freq = i.stationFreq
			}
			if old_name != i.stationName && i.stationName != "\n" {
				data := fmt.Sprintf("name %s\n", i.stationName)
				pipeIn.Write([]byte(data))
				old_name = i.stationName
			}

		}

	}()

	// var stdoutBuf bytes.Buffer

	// w := io.MultiWriter(os.Stdout, &stdoutBuf)
	// cmd_python.Stdout = w
	// go func() {
	// 	fmt.Println(stdoutBuf.String())
	// 	stdoutBuf.Reset()
	// }()

	go func() {
		count := 0
		for {
			var s string
			fmt.Fscan(pipeOut, &s)
			fmt.Println(s, count)
			if s == "freq:-100" {
				i.stationFreq -= 100
			} else if s == "freq:+100" {
				i.stationFreq += 100
			} else {
				fmt.Println(s)
			}

			count += 1
		}
	}()

	//io.Copy(io.Discard, pipeOut)

	cmd_python.Wait()
	return nil
}

func (i *info) run(ctx context.Context) error {

	var err error
	i.stationFreq = 90400 // 96000 //  90900 //  105500 //  104700 //   87800 //

	command_radio := fmt.Sprintf("rtl_fm -M fm -l 0 -A std -p 0 -s 180k -g 30 -F 9 -f %dK", i.stationFreq)
	cmd_radio := exec.CommandContext(ctx, "bash", "-c", command_radio)

	command_audio := "play -v 0.05 -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
	cmd_audio := exec.CommandContext(ctx, "bash", "-c", command_audio)

	// cmd_RDS := exec.CommandContext(ctx, "hexdump", "-C")
	//cmd_RDS := exec.Command("redsea")
	cmd_RDS := exec.CommandContext(ctx, "redsea", "--show-partial", "-r", "180k") // --show-partial

	cmd_python := exec.CommandContext(ctx, "python3", "InputAndScreen.py")

	r_rds, w_rds := io.Pipe()
	r_audio, w_audio := io.Pipe()

	mw := io.MultiWriter(w_rds, w_audio)

	cmd_RDS.Stdin = r_rds
	cmd_audio.Stdin = r_audio
	cmd_audio.Stdout = io.Discard

	cmd_radio.Stdout = mw

	if err != nil {
		return err
	}
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Run audio", cmd_audio.Run())
	}()

	go i.rds(cmd_RDS)
	go i.pythonInterface(cmd_python)
	time.Sleep(1 * time.Second)
	fmt.Println("start radio", cmd_radio.Start())

	fmt.Println("wait radio", cmd_radio.Wait())
	return nil
}

func (i *info) rds(cmd_RDS *exec.Cmd) {

	out, _ := cmd_RDS.StdoutPipe()

	err := cmd_RDS.Start()
	fmt.Println("RDS Started:", err)
	var msg map[string]any
	// var station [8]byte

	dec := json.NewDecoder(out)
	for {
		err := dec.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		//showValue(msg, "ps")
		//showValue(msg, "radiotext")
		if msg["ps"] != nil {
			i.stationName = fmt.Sprintf("%s", msg["ps"])
		} else if msg["partial_ps"] != nil {
			i.stationName = fmt.Sprintf("%s", msg["partial_ps"])
		}

		//showMessage(msg)
		// if msg["ps"] != nil {
		// 	fmt.Println("ps ", msg["ps"])
		// }
		// if msg["partial_ps"] != nil {
		// 	fmt.Println("partial_ps ", msg["partial_ps"])
		// }
		// fmt.Printf("%v\n", msg)
		// fmt.Printf("%v\n", msg["partial_ps"])
		//fmt.Println("-----------------------------------------")
	}

	// io.Copy(os.Stdout, out)

}

func showValue(msg map[string]any, key string) {
	if msg[key] != nil {
		fmt.Println(msg[key])
	} else if msg["partial_"+key] != nil {
		fmt.Println(msg["partial_"+key])
	}
}

func showMessage(msg map[string]any) {
	if msg == nil {
		return
	}
	for k, v := range msg {
		switch v := v.(type) {
		case string:
			fmt.Println(k, ":", v)
		case map[string]any:
			showMessage(v)
		default:
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}
