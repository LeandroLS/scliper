package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

func getFlags() (string, string) {
	var site, links string
	flag.StringVar(&site, "html-from", "", "Site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site to get all links")
	flag.Parse()
	return site, links
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func MakeRequest(site string) *http.Response {
	resp, err := http.Get(site)
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

func main() {
	site, linkSource := getFlags()
	if linkSource != "" {
		GetLinksFrom(linkSource)
	}

	if site != "" {
		DownloadHtmlFromSite(site)
	}

}
