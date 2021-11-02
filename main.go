package main

import (
	"flag"
)

func getFlags() (string, string, string) {
	var site, links, images string
	flag.StringVar(&site, "html-from", "", "Inform a site which you wanna download html")
	flag.StringVar(&links, "links-from", "", "Inform a .html or a site/link to get all links in json file")
	flag.StringVar(&images, "images-from", "", "Inform a .html or a site/link to get all images in a html file")
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
