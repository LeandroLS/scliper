package main

import (
	"io"
	"os"
)

func DownloadHtmlFromSite(site string) {
	resp := MakeRequest(site)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	file := CreateFile(site, ".html")
	WriteInFile(file, body)
	fileCreated := CreateFileStruct(file)
	LogCreatedFileMessage(os.Stdout, fileCreated)
}
