package main 

import (
	"bytes"
	"fmt"
	"golang.org/x/net/ipv4"
	"os/exec"
)

func RunCommand(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if stderr.String() != "" {
		return stderr.String(), err
	}
	return stdout.String(), err

}