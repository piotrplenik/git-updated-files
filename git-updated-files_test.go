package main

import (
	"reflect"
	"testing"
)

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
