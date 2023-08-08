package main

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	//var radio_output bytes.Buffer
	freq := 96000

	command_radio := fmt.Sprintf("rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 30 -F 9 -f %dK", freq)
	cmd_radio := exec.Command("bash", "-c", command_radio)

	r, _ := cmd_radio.StdoutPipe()

	var w io.Writer
	tee := io.TeeReader(r, w)

	wg.Add(1)
	go flux_radio(cmd_radio)

	command_audio := "play -v 0.1 -r 171k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
	cmd_audio := exec.Command("bash", "-c", command_audio)
	cmd_audio.Stdin = r
	wg.Add(2)
	go audio(cmd_audio)

	cmd_RDS := exec.Command("bash", "-c", "redsea --show-partial | jq '.partial_ps'")
	cmd_RDS.Stdin = tee
	wg.Add(3)
	go rds(cmd_RDS)

	for true {
		fmt.Println(cmd_RDS.Stdout)
		if cmd_RDS.Stdout != nil {
			fmt.Println(cmd_RDS.Stdout)
		}

	}

	wg.Wait()

}

func rds(cmd_RDS *exec.Cmd) {
	err := cmd_RDS.Run()
	if err != nil {
		fmt.Print(err)
		return
	}
}

func audio(cmd_audio *exec.Cmd) {
	err := cmd_audio.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func flux_radio(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
