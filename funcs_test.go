package main

import "testing"

func TestCleanString(t *testing.T) {

	checkCorrectMessage := func(t *testing.T, result, expected string) {
		t.Helper()
		if result != expected {
			t.Errorf("result '%s', expected '%s'", result, expected)
		}
	}

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
		checkCorrectMessage(t, result, expected)

	})

	t.Run("Clean String for without protocol site", func(t *testing.T) {
		result := CleanString("www.somehttpsite.com")
		expected := "www.somehttpsite.com"
		checkCorrectMessage(t, result, expected)
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
