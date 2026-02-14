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

func searchSinkInput(sinks []pulseaudio.SinkInput, processes []ps.Process, pid int, action func(pulseaudio.SinkInput) error) error {
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
			return action(input)
		}
	}

	for _, process := range processes {
		if process.PPid() == pid {
			if err := searchSinkInput(sinks, processes, process.Pid(), action); err == nil {
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

	return searchSinkInput(sinks, processes, pid, func(input pulseaudio.SinkInput) error {
		if err := input.ToggleMute(); err != nil {
			return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
		}
		return nil
	})
}

func unmuteSinkInputs(paClient *pulseaudio.Client) error {
	inputs, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	for _, input := range inputs {
		log.Trace().Any("sink_input", input).Msg("unmutting sink input...")

		if err := input.SetMute(false); err != nil {
			return fmt.Errorf("unable to toggle mute of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
		}
	}

	return nil
}

func adjustSinkInputVolume(paClient *pulseaudio.Client, pid int, vc *VolumeChange) error {
	sinks, err := paClient.SinkInputs()
	if err != nil {
		return fmt.Errorf("unable to get pulseaudio sink: %w", err)
	}

	processes, err := ps.Processes()
	if err != nil {
		return err
	}

	return searchSinkInput(sinks, processes, pid, func(input pulseaudio.SinkInput) error {
		currentVolume := input.GetVolume()
		log.Debug().Float32("current_volume", currentVolume).Msg("current volume")

		var newVolume float32
		if vc.IsIncrease {
			newVolume = currentVolume + vc.Amount
		} else {
			newVolume = currentVolume - vc.Amount
			if newVolume < 0 {
				newVolume = 0
			}
		}

		if newVolume > 1.0 {
			newVolume = 1.0
		}

		log.Info().Float32("new_volume", newVolume*100).Msg("setting volume")

		if err := input.SetVolume(newVolume); err != nil {
			return fmt.Errorf("unable to set volume of pulseaudio input '%s'(%d)': %w", input.Name, input.Index, err)
		}

		return nil
	})
}
