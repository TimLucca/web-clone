package main

import (
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

// writeToFile writes a specific URL to a specified file
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

func getStyle(s string) {
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

	// Find and download image URLs
	document.Find("link").Each(func(index int, element *goquery.Selection) {
		lnkSrc, exists := element.Attr("href")
		if exists {
			downloadImage(s, lnkSrc)
		}
	})
}

// findAllImages searches a webpage for all img tags
// it then downloads all of the images
func findAllImages(s string) error {
	// Make HTTP request
	response, err := http.Get(s)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	// Find and download image URLs
	document.Find("img").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("src")
		if exists {
			downloadImage(website, imgSrc)
		}
	})
	return nil
}

// Creates directories needed for files in none exist
func checkDirs(s string) {
	ufile, _ := url.Parse(s)
	seg := strings.Split(ufile.Path, "/")
	last := len(seg) - 1
	if seg[len(seg)-1] == "" {
		last = len(seg) - 2
	}
	dir := s[0 : len(s)-len(seg[last])]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}
}

func downloadImage(w, s string) error {
	f := s
	checkDirs(f)
	if s[:3] != "http" {
		s = w + s
	}
	// f, _ := fileName(s)
	e := writeToFile(s, f)
	if e != nil {
		return e
	}
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

// Finds all children links
func findBranches(s, master string) {
	doc, _ := goquery.NewDocument(s)
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		// ensures not writing same file multiple times
		if _, err := os.Stat(href); os.IsNotExist(err) && (len(href) < 3 || "http" != href[0:4]) {
			link := master + href
			// f, err := fileName(website)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			err := writeToFile(link, href)
			if err != nil {
				log.Println(err)
			}
			if err == nil {
				_ = findAllImages(link)
				findBranches(link, master)
			}
		}

	})
}

// starting place!
func getIndex() {
	// f, err := fileName(website)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(f)

	err := writeToFile(website, "index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = findAllImages(website)
}

func main() {
	getIndex()
	getStyle(website)
	findBranches(website, website)
}
