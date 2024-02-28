package main

import (
	"fmt"
	"os"
	"strconv"

	hypripc "github.com/labi-le/hyprland-ipc-client"
	"github.com/luisnquin/pulseaudio"
)

func main() {
	paClient, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}

	inputs, err := paClient.SinkInputs()
	if err != nil {
		panic(err)
	}

	signature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	hyprClient := hypripc.NewClient(signature)

	window, err := hyprClient.ActiveWindow()
	if err != nil {
		panic(err)
	}

	for _, input := range inputs {
		value := input.PropList["application.process.id"]

		if pid, err := strconv.Atoi(value); err == nil {
			if pid == window.Pid {
				if err := input.ToggleMute(); err != nil {
					panic(err)
				}

				return
			}
		}
	}

	fmt.Println("wuuups...")
}
