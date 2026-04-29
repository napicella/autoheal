package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type composeProject struct {
	// filepath of the yaml file of the project.
	filepath string
	// name of the project.
	name string
}

// Restart the compose project and wait for the restart to complete.
func restartCompose(ctx context.Context, proj composeProject) error {
	// Using the compose CLI instead of the compose go SDK.
	// The compose SDK still relies on the old docker SDK instead of moby.
	// When they switch over to moby we can refactor this to use the compose SDK.
	downCmd := exec.CommandContext(ctx, "docker", "compose", "-f", proj.filepath, "-p", proj.name, "down", "--timeout", "60")
	downCmd.Stdout = os.Stdout
	downCmd.Stderr = os.Stderr
	if err := downCmd.Run(); err != nil {
		return fmt.Errorf("failed to run docker compose down: %w", err)
	}
	upCmd := exec.CommandContext(ctx, "docker", "compose", "-f", proj.filepath, "-p", proj.name, "up", "--remove-orphans", "-d")
	upCmd.Stdout = os.Stdout
	upCmd.Stderr = os.Stderr
	return upCmd.Run()
}
