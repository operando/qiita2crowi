package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/b4b4r07/go-crowi"
)

type Qiita struct {
	Articles []Articles    `json:"articles"`
	Groups   []interface{} `json:"groups"`
	Projects []Projects    `json:"projects"`
	Version  string        `json:"version"`
}

type Articles struct {
	RenderedBody string      `json:"rendered_body"`
	Body         string      `json:"body"`
	Coediting    bool        `json:"coediting"`
	CreatedAt    time.Time   `json:"created_at"`
	Group        interface{} `json:"group"`
	ID           string      `json:"id"`
	Private      bool        `json:"private"`
	Tags         []struct {
		Name     string        `json:"name"`
		Versions []interface{} `json:"versions"`
	} `json:"tags"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	User      struct {
		ID              string `json:"id"`
		PermanentID     int    `json:"permanent_id"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
	Comments []interface{} `json:"comments"`
}

type Projects struct {
	RenderedBody string    `json:"rendered_body"`
	Archived     bool      `json:"archived"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         struct {
		ID              string `json:"id"`
		PermanentID     int    `json:"permanent_id"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
}

var urlSafe = strings.NewReplacer(
	`^`, `＾`, // for Crowi's regexp
	`$`, `＄`,
	`*`, `＊`,
	`%`, `％`, // query
	`?`, `？`,
	`/`, `／`, // Prevent unexpected stratification
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

	client, err := crowi.NewClient(*crowiUrl, *accessToken)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		os.Exit(1)
	}

	for _, article := range q.Articles {
		wg.Add(1)
		go func(a Articles) {
			err := qiita2crowi(client, a)
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

func qiita2crowi(client *crowi.Client, article Articles) error {
	ch <- 1
	defer func() { <-ch }()

	pagePath := getTitlePath(*pagePath, article.Title)
	if !path.IsAbs(pagePath) {
		return fmt.Errorf("%s: invalid page path", pagePath)
	}
	// Create Crowi page
	article.Body = fmt.Sprintf("<!-- Imported by\n%s\n-->\n\n%s",
		strings.TrimLeft(article.URL, "https://"),
		article.Body,
	)
	res, err := client.PagesCreate(pagePath, article.Body)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Error)
	}

	// Download images in the Qiita text
	// then upload to Crowi
	pageId := res.Page.ID
	body := article.Body
	matched := img.FindAllStringSubmatch(article.RenderedBody, -1)
	if len(matched) > 0 {
		for _, urls := range matched {
			for i := 1; i < len(urls); i++ {
				file, err := downloadImage(urls[i])
				if err != nil {
					return err
				}
				res, err := client.AttachmentsAdd(pageId, file)
				if err != nil {
					return err
				}
				if !res.OK {
					return errors.New(res.Error)
				}
				body = strings.Replace(body, urls[i], res.Filename, -1)
			}
		}
		// Update image's links in the Crowi page
		res, err = client.PagesUpdate(pageId, body)
		if err != nil {
			return err
		}
		if !res.OK {
			return errors.New(res.Error)
		}
	}

	// If there are comments, add those at the end of the body
	if len(article.Comments) > 0 {
		body += "# Comments by Qiita:Team\n"
		for _, comment := range article.Comments {
			body += fmt.Sprintf("## %s\n", comment.(map[string]interface{})["user"].(map[string]interface{})["id"].(string))
			body += comment.(map[string]interface{})["body"].(string)
		}
		res, err = client.PagesUpdate(pageId, body)
		if err != nil {
			return err
		}
		if !res.OK {
			return errors.New(res.Error)
		}
	}

	return nil
}

func getApiPath(baseUri, endPoint string) (string, error) {
	base, err := url.Parse(baseUri)
	if err != nil {
		return "", err
	}
	ep, err := url.Parse(endPoint)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(ep).String(), nil
}

func getSaftyPath(path string) string {
	return urlSafe.Replace(path)
}

func getTitlePath(defaultPath, titlePath string) string {
	return path.Clean(path.Join(
		defaultPath,
		getSaftyPath(titlePath),
	))
}

func downloadImage(url string) (string, error) {
	filename := ""
	response, err := http.Get(url)
	if err != nil {
		return filename, err
	}

	// response.Status
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return filename, err
	}

	_, filename = path.Split(url)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return filename, err
	}

	defer func() {
		file.Close()
	}()

	file.Write(body)
	return filename, nil
}
