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

	command := fmt.Sprintf("rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 30 -F 9 -f %dK", freq)
	cmd := exec.Command("bash", "-c", command)
	wg.Add(1)
	go flux_radio(cmd)

	command_audio := "play -v 0.1 -r 171k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
	cmd_audio := exec.Command("bash", "-c", command_audio)

	//récupère le flux de la première commande

	flux, _ := cmd.StdoutPipe()

	//Utilise le Stdout comme entrée pour la commande d'après
	cmd_audio.Stdin = flux
	wg.Add(2)
	go audio(flux, cmd_audio)

	cmd_RDS := exec.Command("bash", "-c", "redsea --show-partial | jq '.partial_ps'")
	cmd_RDS.Stdin =
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

func audio(flux io.ReadCloser, cmd_audio *exec.Cmd) {
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
