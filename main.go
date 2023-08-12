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

func main() {
	// Create a context with cancel function to gracefully handle Ctrl+C events

	ctx, cancel := context.WithCancel(context.Background())
	// Handle Ctrl+C signal (SIGINT)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		fmt.Println("\nCtrl+C received. shutting down...")
		cancel() // Cancel the context when Ctrl+C is received
		os.Exit(1)
	}()

	err := run(ctx)
	if err != nil {
		fmt.Println(err)
	}

}

func pythonInterface(cmd_python *exec.Cmd, freq int, rds) error {
	pipeIn, _ := cmd_python.StdinPipe()
	pipeOut, _ := cmd_python.StdoutPipe()

	err := cmd_python.Start()
	if err != nil {
		return err

	}

	go func() {
		for {
			data := fmt.Sprintf("%d\n", freq)
			pipeIn.Write([]byte(data))
			time.Sleep(time.Second)
		}

	}()
	io.Copy(os.Stdout, pipeOut)
	cmd_python.Wait()
	return nil
}

func run(ctx context.Context) error {

	var err error
	freq := 105500 // 87800 //96000 //90400 //
	command_radio := fmt.Sprintf("rtl_fm -M fm -l 0 -A std -p 0 -s 180k -g 30 -F 9 -f %dK", freq)
	cmd_radio := exec.CommandContext(ctx, "bash", "-c", command_radio)

	command_audio := "play -v 0.05 -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
	cmd_audio := exec.CommandContext(ctx, "bash", "-c", command_audio)

	// cmd_RDS := exec.CommandContext(ctx, "hexdump", "-C")
	//cmd_RDS := exec.Command("redsea")
	cmd_RDS := exec.Command("bash", "-c", "redsea --show-partial -r 180k") // --show-partial

	cmd_python := exec.Command("python3", "InputAndScreen.py")

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

	go rds(cmd_RDS)
	go pythonInterface(cmd_python, freq)
	time.Sleep(1 * time.Second)
	fmt.Println("start radio", cmd_radio.Start())

	fmt.Println("wait radio", cmd_radio.Wait())
	return nil
}

func rds(cmd_RDS *exec.Cmd) {

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

		showValue(msg, "ps")
		showValue(msg, "radiotext")
		showMessage(msg)
		// if msg["ps"] != nil {
		// 	fmt.Println("ps ", msg["ps"])
		// }
		// if msg["partial_ps"] != nil {
		// 	fmt.Println("partial_ps ", msg["partial_ps"])
		// }
		// fmt.Printf("%v\n", msg)
		// fmt.Printf("%v\n", msg["partial_ps"])
		fmt.Println("-----------------------------------------")
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
