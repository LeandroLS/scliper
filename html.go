package main

import (
	"io"
)

func DownloadHtmlFromSite(site string) {
	resp := MakeRequest(site)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	file := CreateFile(site, ".html")
	file.Write(body)
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "HTML")
}
