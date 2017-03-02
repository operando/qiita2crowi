package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	img = regexp.MustCompile(`<img .* ?src="(https?://.*?)"`)
	ch  = make(chan int, 4)
)

var (
	accessToken = flag.String("access-token", "", "Crowi's access token")
	crowiUrl    = flag.String("crowi-url", "", "Your Crowi base URL")
	pagePath    = flag.String("page-path", "/qiita", "Default page path")
)

func main() {
	flag.Parse()

	var q Qiita
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&q); err != nil {
		log.Printf("[ERROR] local json syntax error: %s", err.Error())
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	errs := 0

	for _, article := range q.Articles {
		wg.Add(1)
		go func(a Articles) {
			err := qiita2crowi(a)
			if err != nil {
				log.Printf("[ERROR] %s", err.Error())
				errs++
			}
			wg.Done()
		}(article)
	}
	wg.Wait()

	if errs > 0 {
		log.Printf("Failures %d/%d pages", errs, len(q.Articles))
		os.Exit(1)
	}
}

func qiita2crowi(article Articles) error {
	ch <- 1
	defer func() { <-ch }()

	// Create Crowi page
	article.Body = fmt.Sprintf("<!-- Imported by\n%s\n-->\n\n%s",
		strings.TrimLeft(article.URL, "https://"),
		article.Body,
	)
	crowi, err := crowiPageCreate(article.Title, article.Body)
	if err != nil {
		return err
	}
	if !crowi.OK {
		return fmt.Errorf("Failed to create Crowi page: %s (%s)", article.Title, crowi.Error)
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
					return err
				}
				crowi, err := crowiAttachmentsAdd(pageId, file)
				if err != nil {
					return err
				}
				body = strings.Replace(body, urls[i], crowi.Filename, -1)
			}
		}
		// Update image's links in the Crowi page
		_, err := crowiPageUpdate(pageId, body)
		if err != nil {
			return err
		}
	}

	// If there are comments, add those at the end of the body
	if len(article.Comments) > 0 {
		body += "# Comments by Qiita:Team\n"
		for _, comment := range article.Comments {
			body += fmt.Sprintf("## %s\n", comment.(map[string]interface{})["user"].(map[string]interface{})["id"].(string))
			body += comment.(map[string]interface{})["body"].(string)
		}
		_, err = crowiPageUpdate(pageId, body)
		if err != nil {
			return err
		}
	}

	return nil
}
