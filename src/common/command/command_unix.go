// +build !windows

package command

import (
	"os/exec"
	"syscall"
)

// Kill kill process and group
func (c *Command) Kill() {
	cmd := c.Cmd
	if cmd == nil {
		return
	}
	Kill(cmd)
	_ = cmd.Wait() // avoid error code check
}

// Kill kill
func Kill(cmd *exec.Cmd) {
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		if cmd.SysProcAttr != nil && cmd.SysProcAttr.Setpgid {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
}
