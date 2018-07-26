// +build !windows

package w32

func MessageBox(hwnd HWND, title, caption string, flags uint) int {
	return 0
}
