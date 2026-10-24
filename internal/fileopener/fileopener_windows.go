//go:build windows

package fileopener

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

func Open(path string) error {
	return shellExecute("open", path, "", filepath.Dir(path))
}

func Reveal(path string) error {
	return shellExecute("open", "explorer.exe", fmt.Sprintf(`/select,"%s"`, path), "")
}

func shellExecute(verb, file, parameters, directory string) error {
	verbPtr, err := syscall.UTF16PtrFromString(verb)
	if err != nil {
		return fmt.Errorf("invalid shell verb: %w", err)
	}
	filePtr, err := syscall.UTF16PtrFromString(file)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	parametersPtr, err := syscall.UTF16PtrFromString(parameters)
	if err != nil {
		return fmt.Errorf("invalid shell parameters: %w", err)
	}
	directoryPtr, err := syscall.UTF16PtrFromString(directory)
	if err != nil {
		return fmt.Errorf("invalid working directory: %w", err)
	}

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")
	if err := shell32.Load(); err != nil {
		return fmt.Errorf("cannot load Windows Shell: %w", err)
	}
	result, _, _ := shellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(filePtr)),
		uintptr(unsafe.Pointer(parametersPtr)),
		uintptr(unsafe.Pointer(directoryPtr)),
		1, // SW_SHOWNORMAL
	)
	runtime.KeepAlive(verbPtr)
	runtime.KeepAlive(filePtr)
	runtime.KeepAlive(parametersPtr)
	runtime.KeepAlive(directoryPtr)
	if result <= 32 {
		return fmt.Errorf("Windows could not open the requested item (ShellExecute error %d)", result)
	}
	return nil
}

func configureCommand(*exec.Cmd) {
}
