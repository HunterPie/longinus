package windows

import (
	"fmt"
	"github.com/0xrawsec/golang-win32/win32"
	"github.com/0xrawsec/golang-win32/win32/kernel32"
	"github.com/HunterPie/Longinus/core/reader"
	"golang.org/x/sys/windows"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	allAccess = 0x1F0FFF
)

type Memory struct {
	handle        win32.HANDLE
	processName   string
	buffer        []uint8
	isInitialized bool
}

func (m *Memory) readMemory() {
	moduleHandles, err := kernel32.EnumProcessModules(m.handle)
	if err != nil {
		panic(err)
	}

	var moduleBaseAddress uintptr
	var moduleImageSize uintptr

	for _, moduleHandle := range moduleHandles {
		modulePath, err := kernel32.GetModuleFilenameExW(m.handle, moduleHandle)
		if err != nil {
			panic(err)
		}

		if filepath.Base(modulePath) == m.processName {
			moduleInfo, err := kernel32.GetModuleInformation(m.handle, moduleHandle)
			if err != nil {
				panic(err)
			}
			moduleBaseAddress = uintptr(moduleInfo.LpBaseOfDll)
			moduleImageSize = uintptr(moduleInfo.SizeOfImage)
			break
		}
	}

	if moduleBaseAddress == 0 {
		panic("module not found")
	}

	// Allocate buffer to hold everything we can read
	buffer := make([]byte, moduleImageSize)
	var mbi windows.MemoryBasicInformation
	var offset uintptr

	for offset < moduleImageSize {
		addr := moduleBaseAddress + offset

		// Query region
		if err = windows.VirtualQueryEx(
			windows.Handle(m.handle),
			addr,
			&mbi,
			unsafe.Sizeof(mbi),
		); err != nil {
			panic(fmt.Errorf("VirtualQueryEx failed: %w", err))
		}

		regionSize := uintptr(mbi.RegionSize)
		// clamp to module size
		if offset+regionSize > moduleImageSize {
			regionSize = moduleImageSize - offset
		}

		// Only read committed + accessible regions
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_NOACCESS) == 0 &&
			(mbi.Protect&windows.PAGE_GUARD) == 0 {

			var bytesRead uintptr
			err := windows.ReadProcessMemory(
				windows.Handle(m.handle),
				addr,
				&buffer[offset],
				regionSize,
				&bytesRead,
			)
			if err != nil && err != windows.ERROR_PARTIAL_COPY {
				// we ignore partial copy, but bail on other errors
				panic(fmt.Errorf("ReadProcessMemory failed: %w", err))
			}
		}

		offset += uintptr(mbi.RegionSize) // advance to next region
	}

	m.buffer = buffer
	fmt.Printf("Finished reading module %s into buffer (%d bytes)\n", m.processName, len(buffer))
}

func (m *Memory) Read() []uint8 {
	if !m.isInitialized {
		m.readMemory()
		m.isInitialized = true
	}

	return m.buffer
}

func findProcessIDByName(name string) (uint32, error) {
	// Create snapshot of all processes
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snap)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	// Get first process
	err = windows.Process32First(snap, &entry)
	if err != nil {
		return 0, err
	}

	for {
		// entry.ExeFile is a fixed array [windows.MAX_PATH]uint16
		exe := syscall.UTF16ToString(entry.ExeFile[:])
		if exe == name {
			return entry.ProcessID, nil
		}

		// Next
		err = windows.Process32Next(snap, &entry)
		if err != nil {
			break
		}
	}

	return 0, fmt.Errorf("process %s not found", name)
}

func NewMemory(processName string) reader.ByteDataSource {
	pid, err := findProcessIDByName(processName)
	if err != nil {
		panic(err)
	}
	
	handle, err := kernel32.OpenProcess(allAccess, win32.FALSE, win32.DWORD(pid))
	if err != nil {
		panic(err)
	}
	return &Memory{
		handle:      handle,
		processName: processName,
	}
}
