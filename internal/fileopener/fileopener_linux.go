//go:build linux

package fileopener

import (
	"os/exec"
	"path/filepath"
)

func Open(path string) error {
	return launch("xdg-open", path)
}

func Reveal(path string) error {
	return launch("xdg-open", filepath.Dir(path))
}

func configureCommand(*exec.Cmd) {}
