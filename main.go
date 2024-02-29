package main

import (
	"errors"
	"fmt"
	"os"

	hypripc "github.com/labi-le/hyprland-ipc-client"
	"github.com/luisnquin/go-log"
	"github.com/luisnquin/pulseaudio"
)

func main() {
	signature := os.Getenv(SIGNATURE_ENV_KEY)
	log.Trace().Str("hyprland_is", signature).Send()

	if signature == "" {
		msg := fmt.Sprintf("couldn't get '%s' environment variable, unable to initialize IPC client", SIGNATURE_ENV_KEY)
		log.Error().Msg(msg)
		log.Pretty.Fatal(msg)
	}

	paClient, err := pulseaudio.NewClient()
	if err != nil {
		log.Err(err).Msg("cannot create pulseaudio client, missing pulseaudio or pipewire with 'pipewire-pulse'?")
		log.Pretty.Error1("cannot create pulseaudio client :(")
	}

	hyprClient := hypripc.NewClient(signature)

	window, err := hyprClient.ActiveWindow()
	if err != nil {
		log.Err(err).Msg("couldn't get active Hyprland window...")
		log.Pretty.Error1("couldn't get active Hyprland window")
	}

	if err := toggleSinkInputMute(paClient, window.Pid); err != nil {
		if errors.Is(err, ErrSinkInputNotFound) {
			const msg = "couldn't find a sink input for active window"
			log.Warn().Msg(msg)
			log.Pretty.Error1(msg)
		} else {
			log.Err(err).Msg("couldn't toggle sink input mute...")
			log.Pretty.Error1("couldn't toggle sink input mute :(")
		}
	}
}
