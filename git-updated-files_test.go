package main

import (
	"reflect"
	"testing"
)

func TestGitCurrentBranch(t *testing.T) {
	old := myGitCommand
	defer func() { myGitCommand = old }()

	expect := "feature/test-branch"

	myGitCommand = func(arg ...string) string {
		return expect
	}

	result := gitCurrentBranch()
	if result != expect {
		t.Errorf("Not able to get current branch (result: '%s', expect: '%s')", result, expect)
	}
}

func TestGitGetUpdatedFiles(t *testing.T) {
	old := myGitCommand
	defer func() { myGitCommand = old }()

	expect := "aa\nbb\ncc"

	myGitCommand = func(arg ...string) string {
		return expect
	}

	result := gitGetUpdatedFiles("branchA", "branchB", false)
	if len(result) != 3 {
		t.Errorf("Wrong updated files (result: '%s', expect: '%s')", result, expect)
	}

}

func TestFilterNotMatchedDefaultRule(t *testing.T) {
	files := []string{"file1", "file2", "directory/subdirectory/file3"}
	pattern := ".*"

	result := filterNotMatched(files, pattern)
	expect := files

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("File lists are not equal (result: '%s', expect: '%s')", result, expect)
	}
}

func TestFilterNotMatchedExtensions(t *testing.T) {
	files := []string{"file1.png", "file2.bmp", "directory/subdirectory/file3.png"}
	pattern := ".png"

	result := filterNotMatched(files, pattern)
	expect := []string{"file1.png", "directory/subdirectory/file3.png"}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("File lists are not equal (result: '%s', expect: '%s')", result, expect)
	}
}
