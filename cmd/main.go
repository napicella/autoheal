package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ComposeProject string
	Verbose        bool
	RestartLimit   int
	StopTimeout    int
	CheckInterval  time.Duration
}

func main() {
	cfg := loadConfig()
	initLogger(cfg.Verbose)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create docker client")
	}

	w := NewWatcher(cli, cfg)
	w.Run(ctx)
}

func loadConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.ComposeProject, "project", "", "Only monitor containers belonging to this Docker Compose project")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
	flag.IntVar(&cfg.RestartLimit, "restart-limit", 10, "Maximum number of restarts before stopping container")
	flag.IntVar(&cfg.StopTimeout, "stop-timeout", 10, "Stop timeout (seconds)")
	flag.DurationVar(&cfg.CheckInterval, "interval", 5*time.Second, "Interval between health checks")
	flag.Parse()
	return cfg
}

func initLogger(verbose bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
