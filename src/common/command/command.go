package command

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// Command exec.Cmd wapper
type Command struct {
	*exec.Cmd
}

// coding-staging-k8s/git/tree/master/service/e-git/e-git-rpc-server-ss-1.yml
// Environment variable transfer table
// Add a new environment variable to the whitelist and modify this variable.
var envList = []string{
	"HOME",
	"PATH",
	"LD_LIBRARY_PATH",
	"TZ",
	"AMQP_URL",       // amqpURL
	"GIT_PATH",       // /usr/bin/git
	"DISCOVERY_PATH", // discoveryPath
	"ETCD_MACHINES",  // etcdMachines
	"REDIS_ADDRESS",  // redis host and port
	"REDIS_PASSWORD", // can be empty
	"DEPOT_EVENT_EXPIRED_TIME",
	"AUTH_ADDRESS", // repo-auth-server k8s service name and port
}

// cleanup environ
// Do not use the variable directly!!!
var environ []string
var once sync.Once

func delayInitializeEnv() {
	environ = make([]string, 0, len(envList))
	for _, e := range envList {
		if v := os.Getenv(e); len(v) != 0 {
			environ = append(environ, e+"="+v)
		}
	}
}

// Environ return cleanup environ
func Environ() []string {
	once.Do(delayInitializeEnv)
	return environ
}

// NewNativeCommand native command
func NewNativeCommand(name string, arg ...string) *Command {
	return &Command{Cmd: exec.Command(name, arg...)}
}

// NewNativeCommandContext native command with context
func NewNativeCommandContext(ctx context.Context, name string, arg ...string) *Command {
	return &Command{Cmd: exec.CommandContext(ctx, name, arg...)}
}

// NewCommand todo
func NewCommand(cmd *exec.Cmd, stdin io.Reader, stdout, stderr io.Writer, env ...string) *Command {
	command := &Command{Cmd: cmd}
	delayEnv := Environ()
	cmd.Env = make([]string, 0, len(delayEnv)+len(env))
	cmd.Env = append(cmd.Env, delayEnv...)
	cmd.Env = append(cmd.Env, env...)
	// see: man setpgid
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return command
}

// Pid process pid
func (c *Command) Pid() int {
	cmd := c.Cmd

	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Pid
	}
	return 0
}

// ExitStatus process exit status
func ExitStatus(err error) (int, bool) {
	exitError, ok := err.(*exec.ExitError)
	if !ok {
		return 0, false
	}

	waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
	if !ok {
		return 0, false
	}

	return waitStatus.ExitStatus(), true
}
