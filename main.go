package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	hypr_ipc "github.com/labi-le/hyprland-ipc-client/v3"
	"github.com/luisnquin/go-log"
	"github.com/luisnquin/pulseaudio"
)

const version = "unknown"

func main() {
	debug, showVersion, unmuteAll := false, false, false
	var volumeStr string

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n", APP_NAME)
		flag.PrintDefaults()
	}
	flag.BoolVar(&unmuteAll, "unmute-all", false, "Unmute all pulseaudio sinks(reverts any previous change)")
	flag.StringVar(&volumeStr, "volume", "", "Adjust volume (e.g., '5+', '10-')")
	flag.BoolVar(&showVersion, "version", false, "Print the program version")
	flag.BoolVar(&debug, "debug", false, "Send debug logs to stderr")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	lw, err := getLogsWriter()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer lw.Close()

	if debug {
		log.Init(io.MultiWriter(lw, os.Stderr))
	} else {
		log.Init(lw)
	}

	paClient, err := pulseaudio.NewClient()
	if err != nil {
		log.Err(err).Msg("cannot create pulseaudio client, missing pulseaudio or pipewire with 'pipewire-pulse'?")
		lw.Close()
		log.Pretty.Error1("cannot create pulseaudio client :(")
	}

	var volumeChange *VolumeChange
	if volumeStr != "" {
		volumeChange, err = parseVolumeChange(volumeStr)
		if err != nil {
			log.Err(err).Msg("invalid volume format")
			lw.Close()
			log.Pretty.Error1(fmt.Sprintf("invalid volume format: %s", err.Error()))
		}
		log.Debug().Any("volume_change", volumeChange).Msg("parsed volume change")
	}

	if unmuteAll {
		log.Info().Msg("the goal now is to unmute every sink input")
		if err := unmuteSinkInputs(paClient); err != nil {
			log.Err(err).Msg("cannot unmute sink inputs...")
			lw.Close()
			log.Pretty.Error1(err.Error())
		}
	} else if volumeChange != nil {
		log.Info().Msg("adjust volume of the active window")
		performActiveWindowAction(lw, func(pid int) error {
			return adjustSinkInputVolume(paClient, pid, volumeChange)
		}, "couldn't adjust sink input volume")
	} else {
		log.Info().Msg("mute the sink input of the active window")
		performActiveWindowAction(lw, func(pid int) error {
			return toggleSinkInputMute(paClient, pid)
		}, "couldn't toggle sink input mute")
	}
}

func performActiveWindowAction(lw io.WriteCloser, action func(int) error, errMsg string) {
	signature := os.Getenv(SIGNATURE_ENV_KEY)
	log.Trace().Str("hyprland_is", signature).Send()

	if signature == "" {
		msg := fmt.Sprintf("couldn't get '%s' environment variable, unable to initialize IPC client", SIGNATURE_ENV_KEY)
		log.Error().Msg(msg)
		lw.Close()
		log.Pretty.Fatal(msg)
	}

	hyprClient := hypr_ipc.MustClient(signature)

	window, err := hyprClient.ActiveWindow()
	if err != nil {
		log.Err(err).Msg("couldn't get active Hyprland window...")
		lw.Close()
		log.Pretty.Error1("couldn't get active Hyprland window")
	}

	log.Debug().Any("active_window", window).Send()

	if err := action(window.Pid); err != nil {
		if errors.Is(err, ErrSinkInputNotFound) {
			const msg = "couldn't find a sink input for active window"
			log.Warn().Msg(msg)
			lw.Close()
			log.Pretty.Error1(msg)
		} else {
			log.Err(err).Msgf("%s...", errMsg)
			lw.Close()
			log.Pretty.Error1(fmt.Sprintf("%s :(", errMsg))
		}
	}
}

func getLogsWriter() (io.WriteCloser, error) {
	fp := path.Join(os.TempDir(), fmt.Sprintf("%s.log", APP_NAME))
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND

	f, err := os.OpenFile(fp, flags, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}

	f.WriteString("\n\n")

	return f, nil
}

type VolumeChange struct {
	Amount     float32
	IsIncrease bool
}

func parseVolumeChange(s string) (*VolumeChange, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty volume string")
	}

	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)([\+\-])$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("invalid volume format, expected format: '<number>+' or '<number>-' (e.g., '5+', '10-')")
	}

	amount, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %w", err)
	}

	if amount < 0 {
		return nil, fmt.Errorf("volume amount must be positive")
	}

	vc := &VolumeChange{
		Amount:     float32(amount) / 100.0,
		IsIncrease: matches[2] == "+",
	}

	return vc, nil
}
