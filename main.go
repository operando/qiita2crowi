package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var img = regexp.MustCompile(`<img .* ?src="(.*?)"`)

var (
	debug       = flag.Bool("debug", false, "debug option")
	accessToken = flag.String("access-token", "", "crowi access token")
	crowiUrl    = flag.String("crowi-url", "", "your crowi base url")
	pagePath    = flag.String("page-path", "/qiita", "default path prefix")
)

func main() {
	flag.Parse()
	qiita2crowi(os.Stdin)
}

func qiita2crowi(r io.Reader) {
	dec := json.NewDecoder(r)
	var q Qiita
	dec.Decode(&q)

	for _, article := range q.Articles {
		// Create Crowi page
		crowi, err := crowiPageCreate(article.Title, article.Body)
		if err != nil {
			if *debug {
				log.Printf("[ERROR] %s", err.Error())
			}
			continue
		}
		if !crowi.OK {
			log.Printf("[ERROR] Failed to create Crowi page: %s", article.Title)
			log.Printf(crowi.Error)
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
						if *debug {
							log.Print("[ERROR] %s", err.Error())
						}
						continue
					}
					a, err := crowiAttachmentsAdd(pageId, file)
					if err != nil {
						if *debug {
							log.Print("[ERROR] %s", err.Error())
						}
						continue
					}
					body = strings.Replace(body, urls[i], "/uploads/"+a.Attachment.FilePath, -1)
				}
			}
			// Update image's links in the Crowi page
			_, err := crowiPageUpdate(pageId, body)
			if err != nil {
				if *debug {
					log.Print("[ERROR] %s", err.Error())
				}
				continue
			}
		}
	}
}
