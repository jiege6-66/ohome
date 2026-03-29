//go:build windows

package updater

import "os"

func terminateProcess(process *os.Process) error {
	return process.Kill()
}
