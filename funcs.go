package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	Name      string
	Size      int64
	Extension string
}

func CreateFileStruct(file *os.File) File {
	fileStat, err := file.Stat()
	HandleErr(err)
	fileExtension := filepath.Ext(file.Name())
	fileCreated := File{fileStat.Name(), fileStat.Size(), fileExtension}
	return fileCreated
}

func LogCreatedFileMessage(writer io.Writer, file File) {
	templateStr := `-------------------
File created! 😀
Name: {{ .Name }}
Size: {{ .Size }} bytes
Extension: {{ .Extension }}
-------------------`
	tmpl, err := template.New("LogMessage").Parse(templateStr)
	HandleErr(err)
	err = tmpl.Execute(writer, file)
	HandleErr(err)
}

func GetFlags() (site string, links string, images string) {
	flag.StringVar(&site, "html-from", "", "Inform a site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site/link to get all links in json file")
	flag.StringVar(&images, "images-from", "", "Inform a .html or a site/link to get all images in a html file")
	flag.Parse()
	return site, links, images
}
