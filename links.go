package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
)

func getLinksFromSite(site string) []string {
	resp := MakeRequest(site)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	sReader := strings.NewReader(string(body))
	doc := ParseHtml(sReader)
	HandleErr(err)
	links := GetHtmlTags(doc, "a", "href", nil)
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
	links := GetHtmlTags(doc, "a", "href", nil)
	return links
}

func writeInLinksJsonFile(file *os.File, links []string) {
	var strWithLinks string
	for i := 0; i < len(links); i++ {
		strWithLinks += links[i] + "\n"
	}
	jsonLinks, err := json.MarshalIndent(links, "", "	")
	HandleErr(err)
	fileCreated := CreateFileStruct(file)
	LogCreatedFileMessage(fileCreated, "Links")
	WriteInFile(file, jsonLinks)
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
	file := CreateFile(source, "-links.json")
	writeInLinksJsonFile(file, links)

}
