package main

import (
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

func parseHtml(r io.Reader) *html.Node {
	doc, err := html.Parse(r)
	HandleErr(err)
	return doc
}

func createTxtLinksFile(name string) *os.File {
	name = strings.Trim(name, "https://")
	file, err := os.Create(name + ".txt")
	HandleErr(err)
	return file
}

func getLinksFromSite(site string) []string {
	resp := MakeRequest(site)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := parseHtml(sReader)
	HandleErr(err)
	links := Visit(nil, doc)
	return links
}

func getLinksFromHtml(htmlFile string) []string {
	bytes, err := os.ReadFile(htmlFile)
	HandleErr(err)
	sReader := strings.NewReader(string(bytes))
	doc := parseHtml(sReader)
	links := Visit(nil, doc)
	return links
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
	linksBytes := make([]byte, 0)
	for i := 0; i < len(links); i++ {
		bLink := []byte(links[i])
		linksBytes = append(linksBytes, bLink...)
	}
	file := createTxtLinksFile(source)
	_, err = file.Write(linksBytes)
	HandleErr(err)
}
