package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"
)

const (
	QiitaPathPrefix = "/jp/product/qiita/"
	AccessToken     = ""
)

type Response struct {
	Error string `json:"error"`
	OK    bool   `json:"ok"`
}

type Crowi struct {
	Page struct {
		Revision struct {
			__V       int       `json:"__v"`
			Author    string    `json:"author"`
			Body      string    `json:"body"`
			Path      string    `json:"path"`
			_ID       string    `json:"_id"`
			CreatedAt time.Time `json:"createdAt"`
			Format    string    `json:"format"`
		} `json:"revision"`
		__V            int         `json:"__v"`
		RedirectTo     interface{} `json:"redirectTo"`
		UpdatedAt      time.Time   `json:"updatedAt"`
		LastUpdateUser struct {
			_ID       string    `json:"_id"`
			Name      string    `json:"name"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Status    int       `json:"status"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"lastUpdateUser"`
		Creator struct {
			_ID       string    `json:"_id"`
			Name      string    `json:"name"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Status    int       `json:"status"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"creator"`
		Path         string        `json:"path"`
		_ID          string        `json:"_id"`
		CreatedAt    time.Time     `json:"createdAt"`
		CommentCount int           `json:"commentCount"`
		SeenUsers    []interface{} `json:"seenUsers"`
		Liker        []interface{} `json:"liker"`
		GrantedUsers []string      `json:"grantedUsers"`
		Grant        int           `json:"grant"`
		Status       string        `json:"status"`
		ID           string        `json:"id"`
	} `json:"page"`
	Attachment struct {
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
	} `json:"attachment"`
	Filename string `json:"filename"`
	OK       bool   `json:"ok"`
	Error    string `json:"error"`
}

type Attachments struct {
	Page struct {
		_ID            string      `json:"_id"`
		RedirectTo     interface{} `json:"redirectTo"`
		UpdatedAt      time.Time   `json:"updatedAt"`
		LastUpdateUser struct {
			_ID       string    `json:"_id"`
			Email     string    `json:"email"`
			Username  string    `json:"username"`
			Name      string    `json:"name"`
			Admin     bool      `json:"admin"`
			CreatedAt time.Time `json:"createdAt"`
			Status    int       `json:"status"`
		} `json:"lastUpdateUser"`
		Creator struct {
			_ID       string    `json:"_id"`
			Name      string    `json:"name"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Status    int       `json:"status"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"creator"`
		Path     string `json:"path"`
		__V      int    `json:"__v"`
		Revision struct {
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
		} `json:"revision"`
		CreatedAt    time.Time     `json:"createdAt"`
		CommentCount int           `json:"commentCount"`
		SeenUsers    []string      `json:"seenUsers"`
		Liker        []interface{} `json:"liker"`
		GrantedUsers []string      `json:"grantedUsers"`
		Grant        int           `json:"grant"`
		Status       string        `json:"status"`
		ID           string        `json:"id"`
	} `json:"page"`
	Attachment struct {
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
	} `json:"attachment"`
	Filename string `json:"filename"`
}

func crowiPageCreate(title, body string) (c Crowi, err error) {
	path := QiitaPathPrefix + title

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("access_token", AccessToken)
	writer.WriteField("body", body)
	writer.WriteField("path", path)
	writer.Close()

	resp, err := http.Post(
		"http://localhost:3000/_api/pages.create",
		"multipart/form-data; boundary="+writer.Boundary(),
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

func crowiAttachmentsAdd(pageId, file string) (c Attachments, err error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("access_token", AccessToken)
	writer.WriteField("page_id", pageId)
	{
		header := make(textproto.MIMEHeader)
		header.Add("Content-Disposition", `form-data; name="file"; filename="`+file+`"`)
		header.Add("Content-Type", "image/png")
		fileWriter, err := writer.CreatePart(header)
		if err != nil {
			return c, err
		}
		file, err := os.Open(file)
		if err != nil {
			return c, err
		}
		io.Copy(fileWriter, file)
	}
	writer.Close()

	resp, err := http.Post("http://localhost:3000/_api/attachments.add", "multipart/form-data; boundary="+writer.Boundary(), &buffer)
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

func crowiPageUpdate(pageId, body string) (c Crowi, err error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("access_token", AccessToken)
	writer.WriteField("page_id", pageId)
	writer.WriteField("body", body)
	writer.Close()

	resp, err := http.Post(
		"http://localhost:3000/_api/pages.update",
		"multipart/form-data; boundary="+writer.Boundary(),
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
