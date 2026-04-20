package executil

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CommandContext works like exec.CommandContext, but it ensures that standard 
// macOS paths like /opt/homebrew/bin are searched when looking for the binary.
// This prevents errors like "executable file not found in $PATH" in GUI or background worker apps.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmdName := name
	if path, err := exec.LookPath(name); err == nil {
		cmdName = path
	} else {
		// Look in extended paths
		extendedPaths := []string{
			"/opt/homebrew/bin",
			"/usr/local/bin",
			"/usr/bin",
			"/bin",
		}
		
		// Volta support
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			extendedPaths = append(extendedPaths, filepath.Join(home, ".volta", "bin"))
		}

		for _, dir := range extendedPaths {
			fullPath := filepath.Join(dir, name)
			if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
				if stat.Mode()&0111 != 0 {
					cmdName = fullPath
					break
				}
			}
		}
	}

	cmd := exec.CommandContext(ctx, cmdName, arg...)

	// Also ensure that spawned processes see the extended PATH
	env := os.Environ()
	pathFound := false
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			pathVal := strings.TrimPrefix(e, "PATH=")
			if !strings.Contains(pathVal, "/opt/homebrew/bin") {
				pathVal = "/opt/homebrew/bin:/usr/local/bin:" + pathVal
			}
			env[i] = "PATH=" + pathVal
			pathFound = true
			break
		}
	}
	if !pathFound {
		env = append(env, "PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin")
	}
	cmd.Env = env

	return cmd
}
