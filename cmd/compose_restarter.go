package main

import (
	"context"
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
	restartCmd := exec.CommandContext(ctx, "docker", "compose", "-f", proj.filepath, "-p", proj.name, "restart")
	restartCmd.Stdout = os.Stdout
	restartCmd.Stderr = os.Stderr
	return restartCmd.Run()
}
