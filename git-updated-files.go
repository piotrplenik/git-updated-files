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

func gitCurrentBranch() string {
	gitPath := lookForGitPath()
	cmd := exec.Command(gitPath, "rev-parse", "--abbrev-ref", "HEAD")

	cmdOut, _ := cmd.StdoutPipe()
	cmd.Start()
	bytes, _ := ioutil.ReadAll(cmdOut)

	return strings.TrimSpace(string(bytes))
}

func gitGetUpdatedFiles(sourceBranchName string, targetBranchName string) []string {
	gitPath := lookForGitPath()
	cmd := exec.Command(gitPath, "diff", "--name-only", sourceBranchName, targetBranchName)

	cmdOut, _ := cmd.StdoutPipe()
	cmd.Start()
	bytes, _ := ioutil.ReadAll(cmdOut)

	files := strings.TrimSpace(string(bytes))

	return strings.Split(files, "\n")
}

func validateBranchExist(branchName string) {
	gitPath := lookForGitPath()
	cmd := exec.Command(gitPath, "rev-parse", "--verify", branchName)

	cmd.Start()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 128 {
					log.Printf("Missing branch name '%s'.", branchName)
					syscall.Exit(128)
				}
				log.Fatalf("Error during validaton branch. Status: %d", status.ExitStatus())
				syscall.Exit(1)
			}
		}
		log.Fatalf("Could not receive response: %v", err)
		syscall.Exit(1)
	}
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
