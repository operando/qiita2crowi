package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	PagesCreateAPI    = "/_api/pages.create"
	PagesUpdateAPI    = "/_api/pages.update"
	AttachmentsAddAPI = "/_api/attachments.add"
)

type Crowi struct {
	CrowiPages       CrowiPages
	CrowiAttachments CrowiAttachments
}

type CrowiPages struct {
	Page       PagesPage  `json:"page"`
	Attachment Attachment `json:"attachment"`
	OK         bool       `json:"ok"`
	Error      string     `json:"error"`
}

type CrowiAttachments struct {
	Page       AttachmentsPage `json:"page"`
	Attachment Attachment      `json:"attachment"`
	Filename   string          `json:"filename"`
	OK         bool            `json:"ok"`
	Error      string          `json:"error"`
}

type PagesPage struct {
	Revision       PagesRevision  `json:"revision"`
	_ID            string         `json:"_id"`
	RedirectTo     interface{}    `json:"redirectTo"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	LastUpdateUser LastUpdateUser `json:"lastUpdateUser"`
	Creator        Creator        `json:"creator"`
	Path           string         `json:"path"`
	__V            int            `json:"__v"`
	CreatedAt      time.Time      `json:"createdAt"`
	CommentCount   int            `json:"commentCount"`
	SeenUsers      []string       `json:"seenUsers"`
	Liker          []interface{}  `json:"liker"`
	GrantedUsers   []string       `json:"grantedUsers"`
	Grant          int            `json:"grant"`
	Status         string         `json:"status"`
	ID             string         `json:"id"`
}

type AttachmentsPage struct {
	Revision       AttachmentsRevision `json:"revision"`
	_ID            string              `json:"_id"`
	RedirectTo     interface{}         `json:"redirectTo"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	LastUpdateUser LastUpdateUser      `json:"lastUpdateUser"`
	Creator        Creator             `json:"creator"`
	Path           string              `json:"path"`
	__V            int                 `json:"__v"`
	CreatedAt      time.Time           `json:"createdAt"`
	CommentCount   int                 `json:"commentCount"`
	SeenUsers      []string            `json:"seenUsers"`
	Liker          []interface{}       `json:"liker"`
	GrantedUsers   []string            `json:"grantedUsers"`
	Grant          int                 `json:"grant"`
	Status         string              `json:"status"`
	ID             string              `json:"id"`
}

type Attachment struct {
	__V          int       `json:"__v"`
	FileFormat   string    `json:"fileFormat"`
	FileName     string    `json:"fileName"`
	OriginalName string    `json:"originalName"`
	FilePath     string    `json:"filePath"`
	Creator      string    `json:"creator"`
	Page         string    `json:"page"`
	_ID          string    `json:"_id"`
	CreatedAt    time.Time `json:"createdAt"`
	FileSize     int       `json:"fileSize"`
}

type PagesRevision struct {
	__V       int       `json:"__v"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	_ID       string    `json:"_id"`
	CreatedAt time.Time `json:"createdAt"`
	Format    string    `json:"format"`
}

type AttachmentsRevision struct {
	_ID    string `json:"_id"`
	Author struct {
		_ID       string    `json:"_id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Name      string    `json:"name"`
		Admin     bool      `json:"admin"`
		CreatedAt time.Time `json:"createdAt"`
		Status    int       `json:"status"`
	} `json:"author"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	__V       int       `json:"__v"`
	CreatedAt time.Time `json:"createdAt"`
	Format    string    `json:"format"`
}

type LastUpdateUser struct {
	_ID       string    `json:"_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Admin     bool      `json:"admin"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Creator struct {
	_ID       string    `json:"_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
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
	return (&url.URL{Path: path}).String()
}

func getTitlePath(defaultPath, titlePath string) string {
	return getSaftyPath(path.Clean(path.Join(defaultPath, titlePath)))
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

func crowiPageCreate(title, body string) (c CrowiPages, err error) {
	pagePath := path.Clean(path.Join(*pagePath, title))
	if !path.IsAbs(pagePath) {
		return c, fmt.Errorf("%s: invalid page path", pagePath)
	}

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	w.WriteField("access_token", *accessToken)
	w.WriteField("body", body)
	w.WriteField("path", pagePath)
	w.Close()

	api, err := getApiPath(*crowiUrl, PagesCreateAPI)
	if err != nil {
		return
	}

	resp, err := http.Post(
		api,
		"multipart/form-data; boundary="+w.Boundary(),
		&buffer,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &c)
	if err != nil {
		return
	}
	return
}

func crowiPageUpdate(pageId, body string) (c CrowiPages, err error) {
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	w.WriteField("access_token", *accessToken)
	w.WriteField("page_id", pageId)
	w.WriteField("body", body)
	w.Close()

	api, err := getApiPath(*crowiUrl, PagesUpdateAPI)
	if err != nil {
		return
	}

	resp, err := http.Post(
		api,
		"multipart/form-data; boundary="+w.Boundary(),
		&buffer,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &c)
	if err != nil {
		return
	}
	return
}

func crowiAttachmentsAdd(pageId, file string) (c CrowiAttachments, err error) {
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	w.WriteField("access_token", *accessToken)
	w.WriteField("page_id", pageId)
	{
		header := make(textproto.MIMEHeader)
		header.Add("Content-Disposition", `form-data; name="file"; filename="`+file+`"`)
		header.Add("Content-Type", "image/png")
		fileWriter, err := w.CreatePart(header)
		if err != nil {
			return c, err
		}
		file, err := os.Open(file)
		if err != nil {
			return c, err
		}
		io.Copy(fileWriter, file)
	}
	w.Close()

	api, err := getApiPath(*crowiUrl, AttachmentsAddAPI)
	if err != nil {
		return
	}

	resp, err := http.Post(
		api,
		"multipart/form-data; boundary="+w.Boundary(),
		&buffer,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &c)
	if err != nil {
		return
	}
	return
}
