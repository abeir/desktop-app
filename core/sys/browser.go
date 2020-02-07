package sys

import (
	"errors"
	"fmt"
	"github.com/abeir/desktop-app/core"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error{
	var cmd *exec.Cmd
	switch core.Os() {
	case core.OsMac:
		return errors.New("暂不支持MAC系统调用浏览器")
	case core.OsWindows:
		// 未做测试
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case core.OsLinux:
		cmd = exec.Command("xdg-open", url)
	default:
		return errors.New(fmt.Sprintf("未知的操作系统[%s]，无法调用浏览器", runtime.GOOS))
	}
	return cmd.Start()
}
