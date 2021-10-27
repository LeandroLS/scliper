package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getFlags() (string, string) {
	var siteName, links string
	flag.StringVar(&siteName, "html-from", "", "Site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site to get all links")
	flag.Parse()
	return siteName, links
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createHtmlFile(name string) *os.File {
	name = strings.Trim(name, "https://")
	file, err := os.Create(name + ".html")
	HandleErr(err)
	return file
}

func MakeRequest(siteName string) *http.Response {
	resp, err := http.Get(siteName)
	HandleErr(err)
	return resp
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
		} else {
			return "Links"
		}
	}}).Parse(templateStr)
	HandleErr(err)
	err = tmpl.Execute(os.Stdout, file)
	HandleErr(err)
}

func downloadHtmlFromSite(siteName string) {
	resp := MakeRequest(siteName)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	HandleErr(err)
	file := createHtmlFile(siteName)
	file.Write(body)
	fileStat, err := file.Stat()
	HandleErr(err)
	FileDownloaded := File{fileStat.Name(), fileStat.Size()}
	LogCreatedFileMessage(FileDownloaded, "HTML")
}

func main() {
	siteName, linkSource := getFlags()
	if linkSource != "" {
		GetLinksFrom(linkSource)
	}

	if siteName != "" {
		downloadHtmlFromSite(siteName)
	}

}
