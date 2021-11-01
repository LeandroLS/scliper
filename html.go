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
	file := createHtmlFile(site)
	file.Write(body)
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "HTML")
}

func createHtmlFile(name string) *os.File {
	name = CleanString(name)
	file, err := os.Create(name + ".html")
	HandleErr(err)
	return file
}
