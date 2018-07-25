// +build windows

package w32util

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gabstv/w32"
)

// GetProcessID searches all processes and returns the PID if found
func GetProcessID(name string) (pid uint32, err error) {
	snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPPROCESS, 0)
	if snap == 0 {
		return 0, fmt.Errorf("no snap")
	}
	defer w32.CloseHandle(snap)

	entry := &w32.PROCESSENTRY32{}
	entry.Size = uint32(unsafe.Sizeof(*entry))
	if !w32.Process32First(snap, entry) {
		return 0, fmt.Errorf("not found")
	}
	for w32.Process32Next(snap, entry) {
		n2 := syscall.UTF16ToString(entry.ExeFile[:])
		if name == n2 {
			return entry.ProcessID, nil
		}
	}
	return 0, fmt.Errorf("not found")
}
