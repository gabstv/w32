// +build windows

package w32util

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gabstv/w32"
)

// GetModuleEntryInfo returns the module info (entry point) of a process
func GetModuleEntryInfo(processID uint32, moduleName string) (me *w32.MODULEENTRY32, err error) {
	snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE, processID)
	if snap == 0 {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot returned 0 for the specified PID")
	}
	defer w32.CloseHandle(snap)
	me32 := &w32.MODULEENTRY32{}
	me32.Size = uint32(unsafe.Sizeof(*me32))
	if !w32.Module32First(snap, me32) {
		return nil, fmt.Errorf("not found")
	}
	for hasModule := true; hasModule; hasModule = w32.Module32Next(snap, me32) {
		if syscall.UTF16ToString(me32.SzModule[:]) == moduleName {
			return me32, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

// GetModuleEntries returns all the modules
func GetModuleEntries(processID uint32) ([]w32.MODULEENTRY32, error) {
	modules := make([]w32.MODULEENTRY32, 0)
	snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE, processID)
	if snap == 0 {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot returned 0 for the specified PID")
	}
	defer w32.CloseHandle(snap)
	me32 := &w32.MODULEENTRY32{}
	me32.Size = uint32(unsafe.Sizeof(*me32))
	if !w32.Module32First(snap, me32) {
		return nil, fmt.Errorf("no modules found")
	}
	for hasModule := true; hasModule; hasModule = w32.Module32Next(snap, me32) {
		modules = append(modules, *me32)
	}
	return modules, nil
}
