package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/luisnquin/go-log"
	"github.com/luisnquin/pulseaudio"
	"github.com/mitchellh/go-ps"
)

var ErrSinkInputNotFound = errors.New("sink input couldn't be found")

func searchAndToggleSinkInput(sinks []pulseaudio.SinkInput, processes []ps.Process, pid int) error {
	for _, input := range sinks {
		log.Debug().Any("input", input).Send()

		value := input.PropList[PROCESS_ID_PROPERTY_KEY]

		sinkPid, err := strconv.Atoi(value)
		if err != nil {
			log.Warn().Str("input_name", input.Name).Uint32("input_idx", input.Index).
				Msgf("couldn't convert '%s' to integer, skipping...", value)

			continue
		}

		if sinkPid == pid {
			log.Debug().Msg("sink has been found")

			if err := input.ToggleMute(); err != nil {
				return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
			}

			return nil
		}
	}

	for _, process := range processes {
		if process.PPid() == pid {
			if err := searchAndToggleSinkInput(sinks, processes, process.Pid()); err == nil {
				return nil
			}
		}
	}

	return ErrSinkInputNotFound
}

func toggleSinkInputMute(paClient *pulseaudio.Client, pid int) error {
	sinks, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	processes, err := ps.Processes()
	if err != nil {
		return err
	}

	return searchAndToggleSinkInput(sinks, processes, pid)
}

func unmuteSinkInputs(paClient *pulseaudio.Client) error {
	inputs, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	for _, input := range inputs {
		log.Trace().Any("sink_input", input).Msg("unmutting sink input...")

		if err := input.SetMute(false); err != nil { // collect errors and return everything in a single error set
			return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
		}
	}

	return nil
}
