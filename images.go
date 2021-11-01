package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func VisitImages(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = VisitImages(links, c)
	}

	return links
}

func writeInImgHtmlFile(file *os.File, links []string) {
	var strWithLinks string
	for i := 0; i < len(links); i++ {
		strWithLinks += fmt.Sprintf("<img src=\"%s\"></img> \n", links[i])
	}
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "Images")
	_, err = file.Write([]byte(strWithLinks))
	HandleErr(err)
}

func createHtmlImagesFile(name string) *os.File {
	name = CleanString(name)
	file, err := os.Create(name + "-images.html")
	HandleErr(err)
	return file
}

func GetImages(source string) {
	resp := MakeRequest(source)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := parseHtml(sReader)
	HandleErr(err)
	links := VisitImages(nil, doc)
	links = Map(links, func(link string) string {
		url, err := resp.Request.URL.Parse(link)
		HandleErr(err)
		link = url.String()
		return link
	})
	file := createHtmlImagesFile(source)
	writeInImgHtmlFile(file, links)
}
