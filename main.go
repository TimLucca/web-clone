package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var website = "https://web.archive.org/web/20100528035641im_/http://www.mscd.edu/history/camphale/"

// var file = "abc.html"

func writeToFile(s, file string) error {
	res, err := http.Get(s)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func findAllImages(s string) error {
	// Make HTTP request
	response, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find and print image URLs
	document.Find("img").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("src")
		if exists {
			fmt.Printf("Downloading %s\n", imgSrc)
			downloadImage(s, imgSrc)
		}
	})
	return nil
}

func downloadImage(w, s string) error {
	if s[:3] != "http" {
		s = w + s
	}
	f, _ := fileName(s)
	writeToFile(s, f)
	return nil
}

func fileName(s string) (string, error) {
	ufile, e := url.Parse(s)
	if e != nil {
		return "", e
	}
	seg := strings.Split(ufile.Path, "/")
	name := seg[len(seg)-1]
	if name == "" {
		name = seg[len(seg)-2]
	}
	return name, nil
}

func makeFile(s string) (*os.File, error) {
	file, err := os.Create(s)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func main() {
	f, err := fileName(website)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(f)
	err = writeToFile(website, f)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
	_ = findAllImages(website)
}
