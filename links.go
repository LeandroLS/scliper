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

func parseHtml(r io.Reader) *html.Node {
	doc, err := html.Parse(r)
	HandleErr(err)
	return doc
}

func getLinksFromSite(source string) {
	resp := MakeRequest(source)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := parseHtml(sReader)
	HandleErr(err)
	for _, link := range Visit(nil, doc) {
		fmt.Println(link)
	}
}
func GetLinksFrom(source string) {
	isHtml, err := regexp.MatchString(`\.html$`, source)
	HandleErr(err)
	if isHtml {
		bytes, err := os.ReadFile(source)
		HandleErr(err)
		sReader := strings.NewReader(string(bytes))
		doc := parseHtml(sReader)
		for _, link := range Visit(nil, doc) {
			fmt.Println(link)
		}
	} else {
		getLinksFromSite(source)
	}

}
