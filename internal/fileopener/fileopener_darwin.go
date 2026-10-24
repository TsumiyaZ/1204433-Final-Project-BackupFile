//go:build darwin

package fileopener

import "os/exec"

func Open(path string) error {
	return launch("open", path)
}

func Reveal(path string) error {
	return launch("open", "-R", path)
}

func configureCommand(*exec.Cmd) {}
