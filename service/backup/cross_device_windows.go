//go:build windows

package backup

import (
	"errors"
	"syscall"
)

// errorNotSameDevice is ERROR_NOT_SAME_DEVICE from the Windows API.
const errorNotSameDevice syscall.Errno = 17

func isCrossDeviceError(err error) bool {
	return errors.Is(err, errorNotSameDevice)
}
