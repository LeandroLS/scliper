package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func Visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = Visit(links, c)
	}
	return links
}

func getLinksFromSite(site string) []string {
	resp := MakeRequest(site)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := ParseHtml(sReader)
	HandleErr(err)
	links := Visit(nil, doc)
	links = Map(links, func(link string) string {
		url, err := resp.Request.URL.Parse(link)
		HandleErr(err)
		link = url.String()
		return link
	})
	return links
}

func getLinksFromHtml(htmlFile string) []string {
	bytes, err := os.ReadFile(htmlFile)
	HandleErr(err)
	sReader := strings.NewReader(string(bytes))
	doc := ParseHtml(sReader)
	links := Visit(nil, doc)
	return links
}

func createTxtLinksFile(name string) *os.File {
	name = CleanString(name)
	file, err := os.Create(fmt.Sprintf("%s-links.txt", name))
	HandleErr(err)
	return file
}

func writeInTxtLinksFile(file *os.File, links []string) {
	var strWithLinks string
	for i := 0; i < len(links); i++ {
		strWithLinks += links[i] + "\n"
	}
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "Links")
	_, err = file.Write([]byte(strWithLinks))
	HandleErr(err)
}

func GetLinksFrom(source string) {
	isHtml, err := regexp.MatchString(`\.html$`, source)
	HandleErr(err)
	var links []string
	if isHtml {
		links = append(links, getLinksFromHtml(source)...)
	} else {
		links = append(links, getLinksFromSite(source)...)
	}
	file := createTxtLinksFile(source)
	writeInTxtLinksFile(file, links)

}
