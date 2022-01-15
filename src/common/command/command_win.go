// +build windows

package command

import "os/exec"

// Kill kill command
func (c *Command) Kill() {
	cmd := c.Cmd
	if cmd == nil {
		return
	}
	// cleanup process group
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		_ = cmd.Process.Kill()
	}
	cmd.Wait()
}

// Kill kill
func Kill(cmd *exec.Cmd) {
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		_ = cmd.Process.Kill()
	}
}
