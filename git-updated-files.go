package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
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
	validateGitRepositoryExist()
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

func gitCurrentBranch() string {
	return gitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func gitGetUpdatedFiles(sourceBranchName string, targetBranchName string) []string {
	files := gitCommand("diff", "--name-only", sourceBranchName, targetBranchName)

	return strings.Split(files, "\n")
}

func validateBranchExist(branchName string) {
	if gitCommandStatus("rev-parse", "--verify", branchName) == 128 {
		log.Printf("Missing branch name '%s'.", branchName)
		syscall.Exit(128)
	}
}

func validateGitRepositoryExist() {
	if !gitRepositoryExist() {
		log.Printf("not a git repository (or any of the parent directories): .git")
		syscall.Exit(128)
	}
}

func gitRepositoryExist() bool {
	return gitCommandStatus("status") != 128
}

func filterNotMatched(files []string, pattern string) []string {
	output := []string{}
	for _, file := range files {
		matched, _ := regexp.MatchString(pattern, file)
		if matched {
			output = append(output, file)
		}
	}
	return output
}

func displayFormat(files []string, directoriesOnly bool) []string {
	if !directoriesOnly {
		return files
	}

	keys := make(map[string]bool)

	output := []string{}
	for _, file := range files {
		var entry string

		if directoriesOnly {
			entry = filepath.Dir(file)
		} else {
			entry = file
		}

		if value := keys[entry]; !value && entry != "." {
			keys[entry] = true
			output = append(output, entry)
		}
	}
	return output
}

func main() {
	sourceRefPtr := flag.String("source-ref", gitCurrentBranch(), "GIT source reference")
	targetRefPtr := flag.String("target-ref", "main", "GIT target reference")
	matchPtr := flag.String("filter", ".*", "Regular expression for match")
	directoriesOnlyPtr := flag.Bool("directories-only", true, "Display files only")

	flag.Parse()

	validateGitRepositoryExist()
	validateBranchExist(*sourceRefPtr)
	validateBranchExist(*targetRefPtr)

	files := gitGetUpdatedFiles(*sourceRefPtr, *targetRefPtr)
	filteredFiles := filterNotMatched(files, *matchPtr)

	displayOutput := displayFormat(filteredFiles, *directoriesOnlyPtr)

	stringSlices := strings.TrimSpace(strings.Join(displayOutput[:], "\n"))
	if stringSlices != "" {
		fmt.Println(stringSlices)
	}
	syscall.Exit(0)
}
