//go:build !windows && !darwin && !linux

package fileopener

import (
	"fmt"
	"os/exec"
)

func Open(string) error {
	return fmt.Errorf("opening photos is not supported on this platform")
}

func Reveal(string) error {
	return fmt.Errorf("showing photos in a file manager is not supported on this platform")
}

func configureCommand(*exec.Cmd) {}
