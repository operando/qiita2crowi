package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var img = regexp.MustCompile(`<img .* ?src="(.*?)"`)

var (
	accessToken = flag.String("access-token", "", "crowi access token")
	crowiUrl    = flag.String("crowi-url", "", "your crowi base url")
	crowiPath   = flag.String("crowi-path", "/qiita", "default path prefix")
)

func main() {
	flag.Parse()
	for _, fn := range flag.Args() {
		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		f.Close()
		err = qiita2crowi(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func qiita2crowi(r io.Reader) error {
	dec := json.NewDecoder(r)
	var q Qiita
	dec.Decode(&q)
	for _, article := range q.Articles {
		crowi, err := crowiPageCreate(article.Title, article.Body)
		if err != nil {
			return err
		}
		if !crowi.OK {
			return fmt.Errorf("%s failed to create Crowi page", article.Title)
		}
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
					a, err := crowiAttachmentsAdd(pageId, file)
					if err != nil {
						return err
					}
					body = strings.Replace(body, urls[i], "/uploads/"+a.Attachment.FilePath, -1)
				}
			}
			_, err := crowiPageUpdate(pageId, body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
