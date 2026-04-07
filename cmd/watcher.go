package main

import (
	"context"
	"slices"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

const (
	labelProjName        = "com.docker.compose.project"
	labelComposeFilepath = "com.docker.compose.project.config_files"

	labelAutoheal         = "autoheal"
	labelAutohealStrategy = "autoheal.strategy"
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

	isAutoheal := attrs[labelAutoheal] == "true"
	projName := attrs[labelProjName]
	matchComposeProject := slices.Contains(t.cfg.ComposeProjects, projName)
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

	restartComposeProj := attrs[labelAutohealStrategy] == "project"
	if restartComposeProj {
		log.Info().
			Str("container", name).
			Str("id", ID[:12]).
			Str("project", projName).
			Msg("restarting compose project container")

		composeFilepath := attrs[labelComposeFilepath]
		if composeFilepath == "" {
			log.Error().Msgf(
				"unable to determine compose file path for container, label %s not found",
				labelComposeFilepath)
			return
		}
		log.Info().Msgf("restarting compose project %s", projName)
		err := restartCompose(ctx, composeProject{
			filepath: composeFilepath,
			name:     name,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to restart compose project")
		}
		return
	}

	log.Info().Str("container", name).Str("id", ID[:12]).Msg("restarting unhealthy container")
	timeout := t.cfg.StopTimeout
	if err := t.cli.ContainerRestart(ctx, ID, container.StopOptions{Timeout: &timeout}); err != nil {
		log.Error().Err(err).Str("container", name).Msg("failed to restart container")
	}
}
