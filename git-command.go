package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

func lookForGitPath() string {
	binary, lookErr := exec.LookPath("git")
	if lookErr != nil {
		panic(lookErr)
	}
	return binary
}

func gitCommandStatus(arg ...string) int {
	// validateGitRepositoryExist()
	cmd := exec.Command(lookForGitPath(), arg...)
	cmd.Start()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		log.Fatalf("Could not receive response: %v", err)
		syscall.Exit(1)
	}

	return cmd.ProcessState.ExitCode()
}

func gitCommand(arg ...string) string {
	if !gitRepositoryExist() {
		return ""
	}
	cmd := exec.Command(lookForGitPath(), arg...)

	cmdOut, _ := cmd.StdoutPipe()
	cmd.Start()
	bytes, _ := ioutil.ReadAll(cmdOut)

	return strings.TrimSpace(string(bytes))
}
