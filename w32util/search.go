// +build windows

package w32util

import (
	"fmt"
	"regexp"
	"unsafe"

	"github.com/gabstv/w32"
)

const MEM_IMAGE w32.DWORD = 0x1000000
const MEM_MAPPED w32.DWORD = 0x40000
const MEM_PRIVATE w32.DWORD = 0x20000

// SearchMemory tests a regular expression inside a process memory. It returns
// all the matches.
func SearchMemory(processID, minAddr, maxAddr uint32, maxResults int, exp *regexp.Regexp) ([][]byte, error) {
	phandle, err := w32.OpenProcess(w32.PROCESS_VM_READ|w32.PROCESS_QUERY_INFORMATION, false, processID)
	if err != nil {
		return nil, err
	}
	defer w32.CloseHandle(phandle)
	systemInfo := &w32.SYSTEM_INFO{}
	w32.GetSystemInfo(systemInfo)
	//
	if minAddr == 0 {
		minAddr = uint32(uintptr(systemInfo.MinimumApplicationAddress))
	}
	if maxAddr == 0 {
		maxAddr = uint32(uintptr(systemInfo.MaximumApplicationAddress))
	}
	maxp := uintptr(maxAddr)
	//
	memInfo := &w32.MEMORY_BASIC_INFORMATION{}
	//
	results := make([][]byte, 0)
	//
	for addr := uintptr(minAddr); w32.VirtualQueryEx(phandle, addr, memInfo, int(unsafe.Sizeof(*memInfo))) == int(unsafe.Sizeof(*memInfo)) && addr < maxp; addr += uintptr(memInfo.RegionSize) {
		if memInfo.Protect == w32.PAGE_READWRITE && memInfo.State == w32.MEM_COMMIT {
			buffer, err := w32.ReadProcessMemory(phandle, uint32(addr), uint(memInfo.RegionSize))
			if err != nil {
				fmt.Println("reading error", err.Error())
			} else {
				moreResults := exp.FindAll(buffer, -1)
				if len(moreResults) > 0 {
					results = append(results, moreResults...)
					if maxResults >= 0 && len(results) >= maxResults {
						addr = maxp
						break
					}
				}
			}
		}
	}
	return results, nil
}
