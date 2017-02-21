package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	img     = regexp.MustCompile(`<img .* ?src="(.*?)"`)
	isError = false
)

var (
	debug       = flag.Bool("debug", false, "Logging verbosely")
	accessToken = flag.String("access-token", "", "Crowi's access token")
	crowiUrl    = flag.String("crowi-url", "", "Your Crowi base URL")
	pagePath    = flag.String("page-path", "/qiita", "Default page path")
)

func printError(msg string) {
	if *debug {
		log.Printf("[ERROR] %s", msg)
	}
	isError = true
}

func main() {
	flag.Parse()

	dec := json.NewDecoder(os.Stdin)
	var q Qiita
	dec.Decode(&q)

	for _, article := range q.Articles {
		// Create Crowi page
		crowi, err := crowiPageCreate(article.Title, article.Body)
		if err != nil {
			printError(err.Error())
			continue
		}
		if !crowi.OK {
			printError(fmt.Sprintf("Failed to create Crowi page: %s (%s)", article.Title, crowi.Error))
			continue
		}

		// Download images in the Qiita text
		// then upload to Crowi
		pageId := crowi.Page.ID
		body := article.Body
		matched := img.FindAllStringSubmatch(article.RenderedBody, -1)
		if len(matched) > 0 {
			for _, urls := range matched {
				for i := 1; i < len(urls); i++ {
					file, err := downloadImage(urls[i])
					if err != nil {
						printError(err.Error())
						continue
					}
					a, err := crowiAttachmentsAdd(pageId, file)
					if err != nil {
						printError(err.Error())
						continue
					}
					body = strings.Replace(body, urls[i], "/uploads/"+a.Attachment.FilePath, -1)
				}
			}
			// Update image's links in the Crowi page
			_, err := crowiPageUpdate(pageId, body)
			if err != nil {
				printError(err.Error())
				continue
			}
		}
	}

	if isError {
		os.Exit(1)
	}
}
