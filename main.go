package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type FileDownloaded struct {
	Name string
	Size int64
}

func logCreatedFileMessage(file FileDownloaded) {
	templateStr := `--------------
File: {{ .Name }}
Size: {{ .Size }} bytes
Created
--------------`
	tmpl, err := template.New("test").Parse(templateStr)
	handleErr(err)
	err = tmpl.Execute(os.Stdout, file)
	handleErr(err)
}

func getFlags() (string, string) {
	var siteName, links string
	flag.StringVar(&siteName, "site", "", "Site which you wanna download html")
	flag.StringVar(&links, "links", "", "Inform a .html or a site to get all links")
	flag.Parse()
	return siteName, links
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createHtmlFile(name string) *os.File {
	name = strings.Trim(name, "https://")
	file, err := os.Create(name + ".html")
	handleErr(err)
	return file
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func makeRequest(siteName string) *http.Response {
	resp, err := http.Get(siteName)
	handleErr(err)
	return resp
}

func main() {
	siteName, links := getFlags()

	if links != "" {
		resp := makeRequest(links)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		handleErr(err)
		sReader := strings.NewReader(string(body))
		doc, err := html.Parse(sReader)
		handleErr(err)
		for _, link := range visit(nil, doc) {
			fmt.Println(link)
		}
	}

	if siteName != "" {
		if siteName == "" {
			log.Fatalln("You need to specify a site to download the html")
		}
		resp := makeRequest(siteName)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		handleErr(err)
		file := createHtmlFile(siteName)
		file.Write(body)
		fileStat, err := file.Stat()
		handleErr(err)
		FileDownloaded := FileDownloaded{fileStat.Name(), fileStat.Size()}
		logCreatedFileMessage(FileDownloaded)
	}

}
