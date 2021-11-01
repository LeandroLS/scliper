package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getFlags() (string, string, string) {
	var site, links, images string
	flag.StringVar(&site, "html-from", "", "Inform a site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site/link to get all links")
	flag.StringVar(&images, "images-from", "", "Inform a .html or a site/link to get all images")
	flag.Parse()
	return site, links, images
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
func MakeRequest(site string) *http.Response {
	resp, err := http.Get(site)
	HandleErr(err)
	return resp
}
func ParseHtml(r io.Reader) *html.Node {
	doc, err := html.Parse(r)
	HandleErr(err)
	return doc
}

func CleanString(name string) string {
	name = strings.ReplaceAll(name, "https://", "")
	name = strings.ReplaceAll(name, "/", "-")
	return name
}

type File struct {
	Name string
	Size int64
}

func LogCreatedFileMessage(file File, fileType string) {
	templateStr := `-------------------
{{ type }} file created! 😀
Name: {{ .Name }}
Size: {{ .Size }} bytes
-------------------
`
	tmpl, err := template.New("test").Funcs(template.FuncMap{"type": func() string {
		if fileType == "HTML" {
			return "HTML"
		} else if fileType == "Links" {
			return "Links"
		} else {
			return "Images"
		}
	}}).Parse(templateStr)
	HandleErr(err)
	err = tmpl.Execute(os.Stdout, file)
	HandleErr(err)
}

func main() {
	siteSoure, linkSource, imageSource := getFlags()
	if linkSource != "" {
		GetLinksFrom(linkSource)
	}

	if siteSoure != "" {
		DownloadHtmlFromSite(siteSoure)
	}

	if imageSource != "" {
		GetImages(imageSource)
	}
}
