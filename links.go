package main

import (
	"fmt"
	"io"
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

func GetLinks(links string) {
	resp := MakeRequest(links)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	handleErr(err)
	sReader := strings.NewReader(string(body))
	doc, err := html.Parse(sReader)
	handleErr(err)
	for _, link := range Visit(nil, doc) {
		fmt.Println(link)
	}

}
