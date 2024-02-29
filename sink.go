package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/luisnquin/go-log"
	"github.com/luisnquin/pulseaudio"
)

var ErrSinkInputNotFound = errors.New("sink input couldn't be found")

func toggleSinkInputMute(paClient *pulseaudio.Client, pid int) error {
	inputs, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	for _, input := range inputs {
		value := input.PropList[PROCESS_ID_PROPERTY_KEY]

		sinkPid, err := strconv.Atoi(value)
		if err != nil {
			log.Warn().Str("input_name", input.Name).Uint32("input_idx", input.Index).
				Msgf("couldn't convert '%s' to integer, skipping...", value)

			continue
		}

		if sinkPid == pid {
			if err := input.ToggleMute(); err != nil {
				return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
			}

			return nil
		}
	}

	return ErrSinkInputNotFound
}

func unmuteSinkInputs(paClient *pulseaudio.Client) error {
	inputs, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	for _, input := range inputs {
		if err := input.SetMute(false); err != nil { // collect errors and return everything in a single error set
			return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
		}
	}

	return nil
}
