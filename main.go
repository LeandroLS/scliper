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

type FileDownloaded struct {
	Name string
}

func logCreatedFileMessage(file FileDownloaded) {
	templateStr := "File {{ .Name }} created"
	tmpl, err := template.New("test").Parse(templateStr)
	handleErr(err)
	err = tmpl.Execute(os.Stdout, file)
	handleErr(err)
}

func getSiteName() string {
	var siteName string
	flag.StringVar(&siteName, "site", "", "Site which you wanna download html")
	flag.Parse()
	return siteName
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

func main() {
	siteName := getSiteName()
	if siteName == "" {
		log.Fatalln("You need to specify a site to download the html")
	}
	resp, err := http.Get(siteName)
	handleErr(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	handleErr(err)
	file := createHtmlFile(siteName)
	file.Write(body)
	FileDownloaded := FileDownloaded{file.Name()}
	logCreatedFileMessage(FileDownloaded)
}
