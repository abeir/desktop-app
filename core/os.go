package core

import (
	"runtime"
)

type OsType string

const (
	OsWindows OsType  = "windows"
	OsLinux OsType 	  = "linux"
	OsMac OsType	  = "darwin"
	OsUnknow OsType  = "unknow"
)

type ArchType string
const (
	ArchX86 ArchType 		= "386"
	ArchX64 ArchType 		= "amd64"
	ArchUnknow ArchType 	= "unknow"
)

func Os() OsType{
	switch runtime.GOOS {
	case "darwin":
		return OsMac
	case "windows":
		return OsWindows
	case "linux":
		return OsLinux
	default:
		return OsUnknow
	}
}

func Arch() ArchType {
	switch runtime.GOARCH {
	case "386":
		return ArchX86
	case "amd64":
		return ArchX64
	default:
		return ArchUnknow
	}
}