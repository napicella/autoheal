package main

import (
	"context"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

type Watcher struct {
	cli              *client.Client
	cfg              Config
	numberOfRestarts map[string]int
}

func NewWatcher(cli *client.Client, cfg Config) *Watcher {
	return &Watcher{
		numberOfRestarts: make(map[string]int),
		cfg:              cfg,
		cli:              cli,
	}
}

func (t *Watcher) Run(ctx context.Context) {
	args := filters.NewArgs()
	args.Add("event", "health_status: unhealthy")

	msgs, errs := t.cli.Events(ctx, events.ListOptions{
		Filters: args,
		Since:   (1 * time.Second).String(),
	})

	for {
		select {
		case <-ctx.Done():
			log.Warn().Msg("Closing the watcher")
			return
		case err := <-errs:
			log.Warn().Err(err).Msg("event stream error")
			time.Sleep(2 * time.Second)
		case msg := <-msgs:
			name := msg.Actor.Attributes["name"]

			log.Debug().
				Str("container", name).
				Str("event", string(msg.Action)).
				Msg("Container event")

			t.restart(ctx, msg)
		}
	}
}

func (t *Watcher) restart(ctx context.Context, msg events.Message) {
	attrs := msg.Actor.Attributes
	ID := msg.Actor.ID

	isAutoheal := attrs["autoheal"] == "true"
	matchComposeProject := true
	if t.cfg.ComposeProject != "" {
		matchComposeProject = attrs["com.docker.compose.project"] == t.cfg.ComposeProject
	}
	shouldRestart := isAutoheal && matchComposeProject
	if !shouldRestart {
		return
	}
	name := attrs["name"]
	t.numberOfRestarts[name]++
	restarts := t.numberOfRestarts[name]

	if restarts >= t.cfg.RestartLimit {
		log.Warn().Str("container", name).Msg("max restarts reached, stopping container")
		timeout := t.cfg.StopTimeout
		if err := t.cli.ContainerStop(ctx, ID, container.StopOptions{Timeout: &timeout}); err != nil {
			log.Error().Err(err).Str("container", name).Msg("failed to stop container")
		}
		return
	}

	log.Info().Str("container", name).Str("id", ID[:12]).Msg("restarting unhealthy container")
	timeout := t.cfg.StopTimeout
	if err := t.cli.ContainerRestart(ctx, ID, container.StopOptions{Timeout: &timeout}); err != nil {
		log.Error().Err(err).Str("container", name).Msg("failed to restart container")
	}
}
