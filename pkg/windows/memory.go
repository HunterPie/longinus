package windows

import (
	"github.com/0xrawsec/golang-win32/win32"
	"github.com/0xrawsec/golang-win32/win32/kernel32"
	"golang.org/x/sys/windows"
	"path/filepath"
	"unsafe"
)

type Memory struct {
	handle        win32.HANDLE
	processName   string
	buffer        []uint8
	isInitialized bool
}

func (m *Memory) readMemory() {
	moduleHandles, _ := kernel32.EnumProcessModules(m.handle)
	var moduleBaseAddress uintptr
	var moduleImageSize uintptr

	for _, moduleHandle := range moduleHandles {
		modulePath, _ := kernel32.GetModuleFilenameExW(m.handle, moduleHandle)

		if filepath.Base(modulePath) == m.processName {
			moduleInfo, _ := kernel32.GetModuleInformation(m.handle, moduleHandle)
			moduleBaseAddress = uintptr(moduleInfo.LpBaseOfDll)
			moduleImageSize = uintptr(moduleInfo.SizeOfImage)
			break
		}
	}
	var lpBytesRead uintptr
	bytesBuffer := make([]byte, moduleBaseAddress+moduleImageSize)
	unsafeBytesPtr := unsafe.Pointer(&bytesBuffer[0])

	_ = windows.ReadProcessMemory(
		windows.Handle(m.handle),
		moduleBaseAddress,
		(*byte)(unsafeBytesPtr),
		moduleBaseAddress+moduleImageSize,
		&lpBytesRead,
	)

	m.buffer = bytesBuffer
}

func (m *Memory) Read() []uint8 {
	if !m.isInitialized {
		m.readMemory()
		m.isInitialized = true
	}

	return m.buffer
}
