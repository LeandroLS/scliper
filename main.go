package main

func main() {
	siteSoure, linkSource, imageSource := GetFlags()
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
