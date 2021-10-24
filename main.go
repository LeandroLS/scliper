package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createHtmlFile(name string) *os.File {
	file, err := os.Create(name + ".html")
	handleErr(err)
	return file
}

func main() {
	resp, err := http.Get("https://google.com.br")
	handleErr(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	handleErr(err)
	file := createHtmlFile("teste")
	file.Write(body)
}
