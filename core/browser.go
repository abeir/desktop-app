package core

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error{
	var cmd *exec.Cmd
	switch Os() {
	case OsMac:
		return errors.New("暂不支持MAC系统调用浏览器")
	case OsWindows:
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case OsLinux:
		cmd = exec.Command("xdg-open", url)
	default:
		return errors.New(fmt.Sprintf("未知的操作系统[%s]，无法调用浏览器", runtime.GOOS))
	}
	return cmd.Start()
}
