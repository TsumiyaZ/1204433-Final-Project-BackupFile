package fileopener

import (
	"fmt"
	"os/exec"
)

func launch(name string, args ...string) error {
	command := exec.Command(name, args...)
	configureCommand(command)
	if err := command.Start(); err != nil {
		return fmt.Errorf("cannot start %s: %w", name, err)
	}
	go func() { _ = command.Wait() }()
	return nil
}
