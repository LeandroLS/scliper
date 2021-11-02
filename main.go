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

func CreateFile(name string, suffix string) *os.File {
	name = CleanString(name)
	file, err := os.Create(fmt.Sprintf("%s%s", name, suffix))
	HandleErr(err)
	return file
}

func WriteInFile(file *os.File, content []byte) {
	_, err := file.Write(content)
	HandleErr(err)
}

func GetHtmlTags(n *html.Node, rawHtmlTag string, htmlTagKey string, tags []string) []string {
	if n.Type == html.ElementNode && n.Data == rawHtmlTag {
		for _, a := range n.Attr {
			if a.Key == htmlTagKey {
				tags = append(tags, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tags = GetHtmlTags(c, rawHtmlTag, htmlTagKey, tags)
	}
	return tags
}

func CleanString(name string) string {
	name = strings.ReplaceAll(name, "https://", "")
	name = strings.ReplaceAll(name, "http://", "")
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
			return "HTML .html"
		} else if fileType == "Links" {
			return "Links .json"
		} else {
			return "Images .html"
		}
	}}).Parse(templateStr)
	HandleErr(err)
	err = tmpl.Execute(os.Stdout, file)
	HandleErr(err)
}

func getFlags() (string, string, string) {
	var site, links, images string
	flag.StringVar(&site, "html-from", "", "Inform a site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site/link to get all links")
	flag.StringVar(&images, "images-from", "", "Inform a .html or a site/link to get all images")
	flag.Parse()
	return site, links, images
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
