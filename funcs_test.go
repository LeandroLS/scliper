package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func checkIfIsExpected(t *testing.T, result, expected string) {
	t.Helper()
	if result != expected {
		t.Errorf("result '%s', expected '%s'", result, expected)
	}
}

func TestCleanString(t *testing.T) {

	t.Run("Clean String for https site", func(t *testing.T) {
		result := CleanString("https://somehttpssite.com")
		expected := "somehttpssite.com"
		if result != expected {
			t.Errorf("result '%s', expected '%s'", result, expected)
		}
	})

	t.Run("Clean String for http site", func(t *testing.T) {
		result := CleanString("http://somehttpsite.com")
		expected := "somehttpsite.com"
		checkIfIsExpected(t, result, expected)

	})

	t.Run("Clean String for without protocol site", func(t *testing.T) {
		result := CleanString("www.somehttpsite.com")
		expected := "www.somehttpsite.com"
		checkIfIsExpected(t, result, expected)
	})
}

func TestCreateFile(t *testing.T) {
	tmpdir := t.TempDir()
	file := CreateFile(tmpdir+"TestCreateFile", ".html")
	fileStat, _ := file.Stat()
	expected := fileStat.Name()
	if fileStat.Name() != expected {
		t.Errorf("result '%s', expected '%s'", fileStat.Name(), expected)
	}
	file.Close()
}

func TestGetFlags(t *testing.T) {
	links, site, images := GetFlags()
	if links != "" || site != "" || images != "" {
		t.Errorf("Flags is not empty. Flags: %s, %s, %s", links, site, images)
	}
}

func TestCreateFileStruct(t *testing.T) {
	tmpdir := t.TempDir()
	extension := ".html"
	file := CreateFile(tmpdir+"TestCreateFile", extension)
	fileStruct := CreateFileStruct(file)
	file.Close()
	checkIfIsExpected(t, fileStruct.Extension, extension)
}

func TestGetHtmlTags(t *testing.T) {
	sReader := strings.NewReader(`
	<a href="test1.golang">
	<a href="test2.golang">
`)
	htmlNode := ParseHtml(sReader)
	tags := GetHtmlTags(htmlNode, "a", "href", nil)
	expected := 2
	if len(tags) < expected {
		t.Errorf("result '%d', expected '%d'", len(tags), expected)
	}
}

func TestLogCreatedFileMessage(t *testing.T) {
	extension := ".html"
	tmpdir := t.TempDir()
	file := CreateFile(tmpdir+"TestCreateFile", extension)
	fileStruct := CreateFileStruct(file)
	fileStat, _ := file.Stat()
	file.Close()
	buffer := bytes.Buffer{}
	LogCreatedFileMessage(&buffer, fileStruct)
	expected := fmt.Sprintf(`-------------------
File created! 😀
Name: %s
Size: 0 bytes
Extension: .html
-------------------`, fileStat.Name())
	checkIfIsExpected(t, buffer.String(), expected)
}
