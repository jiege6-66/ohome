//go:build !windows

package updater

import (
	"os"
	"syscall"
)

func terminateProcess(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
