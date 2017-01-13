// Based on ssh/terminal:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows,!appengine

package logrus

import (
	"fmt"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
)

// IsTerminal returns true if stderr's file descriptor is a terminal.
func IsTerminal() bool {
	fd := syscall.Stderr
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

const (
	INVALID_HANDLE_VALUE = ^uintptr(0)

	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	STD_OUTPUT_HANDLE                  = uint32(-11 & 0xFFFFFFFF)
)

func init() {
	err := EnableVTMode()
	if err != nil {
		panic(err)
	}
}

// Console Virtual Terminal Sequences
// https://msdn.microsoft.com/en-us/library/windows/desktop/mt638032(v=vs.85).aspx
func EnableVTMode() error {
	dll := syscall.MustLoadDLL("kernel32.dll")
	GetStdHandle := dll.MustFindProc("GetStdHandle")
	GetConsoleMode := dll.MustFindProc("GetConsoleMode")
	SetConsoleMode := dll.MustFindProc("SetConsoleMode")

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	// HANDLE WINAPI GetStdHandle(
	r, _, err := GetStdHandle.Call(uintptr(STD_OUTPUT_HANDLE))
	fmt.Printf("%T, %d\n", err, err)
	if r == uintptr(INVALID_HANDLE_VALUE) || r == 0 {
		return err
	}
	hout := syscall.Handle(r)
	var dwmode uint32 // DWORD
	// BOOL WINAPI GetConsoleMode(
	r, _, err = GetConsoleMode.Call(uintptr(hout), uintptr(unsafe.Pointer(&dwmode)))
	if r == 0 {
		return err
	}
	//fmt.Println(dwmode | ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	// BOOL WINAPI SetConsoleMode(
	r, _, err = SetConsoleMode.Call(uintptr(hout), uintptr(dwmode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
	if r == 0 {
		return err
	}
	return nil
}
