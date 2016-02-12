// need to change to package api
package main

import (
	"fmt"
	"os/exec"
)

type SSHCommander struct {
	User string
	IP   string
}

func (s *SSHCommander) Command(cmd ...string) *exec.Cmd {
	arg := append(
		[]string{
			fmt.Sprintf("%s@%s", s.User, s.IP),
		},
		cmd...,
	)
	return exec.Command("ssh", arg...)
}
