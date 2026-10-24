//go:build !windows

package llama

import "os/exec"

func configureServerProcess(command *exec.Cmd) {}
