package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAutohealE2E(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	composeFile := "./testdata/docker-compose.test.yml"
	testdataDir := t.TempDir()

	t.Logf("Using temp dir for bind mounts: %s", testdataDir)

	require.NoError(t, os.MkdirAll(filepath.Join(testdataDir, "should-keep-restarting"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(testdataDir, "shouldnt-restart-healthy"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(testdataDir, "shouldnt-restart-no-label"), 0755))

	// Start environment
	up := exec.Command("docker", "compose", "-p", "teststack", "-f", composeFile, "up", "-d", "--build")
	up.Env = append(os.Environ(), fmt.Sprintf("TESTDATA_DIR=%s", testdataDir))
	up.Stdout, up.Stderr = os.Stdout, os.Stderr
	require.NoError(t, up.Run())

	defer func() {
		down := exec.Command("docker", "compose", "-p", "teststack", "-f", composeFile, "down", "-v", "--remove-orphans")
		down.Env = append(os.Environ(), fmt.Sprintf("TESTDATA_DIR=%s", testdataDir))
		down.Stdout, down.Stderr = os.Stdout, os.Stderr
		_ = down.Run()
	}()

	t.Log("Waiting for autoheal to act, then stop it")
	time.Sleep(35 * time.Second)
	// cancel()

	// Count restarts by number of lines in log
	countLines := func(path string) int {
		data, err := os.ReadFile(path)
		if err != nil {
			return 0
		}
		return len(strings.Split(strings.TrimSpace(string(data)), "\n"))
	}

	unhealthyRestarts := countLines(filepath.Join(testdataDir, "should-keep-restarting", "start.log"))
	healthyRestarts := countLines(filepath.Join(testdataDir, "shouldnt-restart-healthy", "start.log"))
	noLabelRestarts := countLines(filepath.Join(testdataDir, "shouldnt-restart-no-label", "start.log"))

	t.Logf("should-keep-restarting=%d, shouldnt-restart-healthy=%d, shouldnt-restart-no-label=%d",
		unhealthyRestarts, healthyRestarts, noLabelRestarts)

	require.GreaterOrEqual(t, unhealthyRestarts, 2, "expected unhealthy container to be restarted at least once")
	require.Equal(t, 1, healthyRestarts, "healthy container should not restart")
	require.Equal(t, 1, noLabelRestarts, "no-label container should not restart")
}
