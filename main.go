package main

import (
	"fmt"
	"os"

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

	for _, input := range inputs {
		fmt.Println(input)
	}

	hyprClient := hypripc.NewClient(os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"))
	fmt.Println("hyprClient.ActiveWindow()")
	fmt.Println(hyprClient.ActiveWindow())
}

// pactl list sink-inputs | grep --before-context=30 --after-context=100 spotify
// pactl set-sink-input-mute 198 toggle
