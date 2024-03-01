package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	hypr_ipc "github.com/labi-le/hyprland-ipc-client"
	"github.com/luisnquin/go-log"
	"github.com/luisnquin/pulseaudio"
)

var version = "unknown"

func main() {
	debug, showVersion, unmuteAll := false, false, false

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n", APP_NAME)
		flag.PrintDefaults()
	}
	flag.BoolVar(&unmuteAll, "unmute-all", false, "Unmute all pulseaudio sinks(reverts any previous change)")
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
	log.Init(io.MultiWriter(lw, os.Stderr))

	paClient, err := pulseaudio.NewClient()
	if err != nil {
		log.Err(err).Msg("cannot create pulseaudio client, missing pulseaudio or pipewire with 'pipewire-pulse'?")
		lw.Close()
		log.Pretty.Error1("cannot create pulseaudio client :(")
	}

	if unmuteAll {
		if err := unmuteSinkInputs(paClient); err != nil {
			log.Err(err).Msg("cannot unmute sink inputs...")
			lw.Close()
			log.Pretty.Error1(err.Error())
		}
	} else {
		signature := os.Getenv(SIGNATURE_ENV_KEY)
		log.Trace().Str("hyprland_is", signature).Send()

		if signature == "" {
			msg := fmt.Sprintf("couldn't get '%s' environment variable, unable to initialize IPC client", SIGNATURE_ENV_KEY)
			log.Error().Msg(msg)
			lw.Close()
			log.Pretty.Fatal(msg)
		}

		hyprClient := hypr_ipc.NewClient(signature)

		window, err := hyprClient.ActiveWindow()
		if err != nil {
			log.Err(err).Msg("couldn't get active Hyprland window...")
			lw.Close()
			log.Pretty.Error1("couldn't get active Hyprland window")
		}

		log.Debug().Any("active_window", window).Send()

		if err := toggleSinkInputMute(paClient, window.Pid); err != nil {
			if errors.Is(err, ErrSinkInputNotFound) {
				const msg = "couldn't find a sink input for active window"
				log.Warn().Msg(msg)
				lw.Close()
				log.Pretty.Error1(msg)
			} else {
				log.Err(err).Msg("couldn't toggle sink input mute...")
				lw.Close()
				log.Pretty.Error1("couldn't toggle sink input mute :(")
			}
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
