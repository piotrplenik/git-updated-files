package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// required for mock testing
var myGitCommand = gitCommand

func gitCurrentBranch() string {
	return myGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func gitGetUpdatedFiles(sourceBranchName string, targetBranchName string, deletedOnly bool) []string {
	// By default filter only `A`-Added, `C`-Copied, `M`-Modified, `R`-Renamed
	diffFilter := "ACMR"
	if deletedOnly {
		// When deletedOnly is enabled, filter `D`-Deleted only
		diffFilter = "D"
	}
	files := myGitCommand("diff", "--name-only", "--diff-filter="+diffFilter, sourceBranchName, targetBranchName)

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
	deletedOnlyPtr := flag.Bool("deleted", false, "Display deleted only")

	flag.Parse()

	validateGitRepositoryExist()
	validateBranchExist(*sourceRefPtr)
	validateBranchExist(*targetRefPtr)

	files := gitGetUpdatedFiles(*sourceRefPtr, *targetRefPtr, *deletedOnlyPtr)
	filteredFiles := filterNotMatched(files, *matchPtr)

	displayOutput := displayFormat(filteredFiles, *directoriesOnlyPtr)

	stringSlices := strings.TrimSpace(strings.Join(displayOutput[:], "\n"))
	if stringSlices != "" {
		fmt.Println(stringSlices)
	}
	syscall.Exit(0)
}
