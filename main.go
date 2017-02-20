package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
)

var img = regexp.MustCompile(`<img .* ?src="(.*?)"`)

func main() {
	dec := json.NewDecoder(os.Stdin)
	var q Qiita
	dec.Decode(&q)
	for _, article := range q.Articles {
		crowi, err := crowiPageCreate(article.Title, article.Body)
		if err != nil {
			panic(err)
		}
		if !crowi.OK {
			log.Printf("%s failed to create Crowi page", article.Title)
		}
		pageId := crowi.Page.ID
		body := article.Body
		matched := img.FindAllStringSubmatch(article.RenderedBody, -1)
		if len(matched) > 0 {
			for _, urls := range matched {
				for i := 1; i < len(urls); i++ {
					file, err := downloadImage(urls[i])
					if err != nil {
						panic(err)
					}
					a, err := crowiAttachmentsAdd(pageId, file)
					if err != nil {
						panic(err)
					}
					body = strings.Replace(body, urls[i], "/uploads/"+a.Attachment.FilePath, -1)
				}
			}
			_, err := crowiPageUpdate(pageId, body)
			if err != nil {
				panic(err)
			}
		}
	}
}
