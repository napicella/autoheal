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
	ComposeProjects projectsFlag
	Verbose         bool
	RestartLimit    int
	StopTimeout     int
	CheckInterval   time.Duration
}

func main() {
	cfg := loadConfig()
	initLogger(cfg.Verbose)
	log.Info().Msgf("Autoheal config: %+v", cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create docker client")
	}

	log.Info().Msgf("Monitoring unhealthy containers labelled 'autoheal' from projects %s", cfg.ComposeProjects)
	w := NewWatcher(cli, cfg)
	w.Run(ctx)
}

func loadConfig() Config {
	var cfg Config
	flag.Var(&cfg.ComposeProjects, "project",
		"Comma-separated list of Docker Compose project to monitor. Only monitor containers belonging to the projects listed")
	flag.BoolVar(&cfg.Verbose, "verbose", false,
		"Enable verbose logging")
	flag.IntVar(&cfg.RestartLimit, "restart-limit", 10,
		"Maximum number of restarts before stopping container")
	flag.IntVar(&cfg.StopTimeout, "stop-timeout", 10,
		"Stop timeout (seconds)")
	flag.DurationVar(&cfg.CheckInterval, "interval", 5*time.Second,
		"Interval between health checks")
	flag.Parse()
	return cfg
}

func initLogger(verbose bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	printBanner()

	logLevel := zerolog.InfoLevel
	if verbose {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
}
