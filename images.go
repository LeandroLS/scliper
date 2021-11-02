package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func writeInImgHtmlFile(file *os.File, links []string) {
	var strWithLinks string
	for i := 0; i < len(links); i++ {
		strWithLinks += fmt.Sprintf("<img src=\"%s\"></img> \n", links[i])
	}
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "Images")
	WriteInFile(file, []byte(strWithLinks))
}

func GetImages(source string) {
	resp := MakeRequest(source)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := ParseHtml(sReader)
	HandleErr(err)
	links := GetHtmlTags(doc, "img", "href", nil)
	links = Map(links, func(link string) string {
		url, err := resp.Request.URL.Parse(link)
		HandleErr(err)
		link = url.String()
		return link
	})
	file := CreateFile(source, "-images.html")
	writeInImgHtmlFile(file, links)
}
