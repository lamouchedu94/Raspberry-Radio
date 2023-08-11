package main

import (
	"context"
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
		fmt.Println("\nCtrl+C received. Gracefully shutting down...")
		cancel() // Cancel the context when Ctrl+C is received
		os.Exit(1)
	}()

	err := run(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func run(ctx context.Context) error {

	var err error
	freq := 90400
	command_radio := fmt.Sprintf("rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 30 -F 9 -f %dK", freq)
	cmd_radio := exec.CommandContext(ctx, "bash", "-c", command_radio)

	command_audio := "play -v 0.05 -r 171k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
	cmd_audio := exec.CommandContext(ctx, "bash", "-c", command_audio)

	// cmd_RDS := exec.CommandContext(ctx, "hexdump", "-C")
	//cmd_RDS := exec.Command("redsea")
	cmd_RDS := exec.Command("bash", "-c", "redsea --show-partial | grep partial_ps")

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

	time.Sleep(1 * time.Second)
	fmt.Println("start radio", cmd_radio.Start())

	fmt.Println("wait radio", cmd_radio.Wait())
	return nil
}

func rds(cmd_RDS *exec.Cmd) {

	out, _ := cmd_RDS.StdoutPipe()

	err := cmd_RDS.Start()
	fmt.Println("RDS Started:", err)

	io.Copy(os.Stdout, out)

}
