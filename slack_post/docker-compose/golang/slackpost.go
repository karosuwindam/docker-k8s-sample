package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SlackApi struct {
	Token   string
	Channel string
}

const (
	CHANNEL          = "general"
	SLACKAPIURL      = "https://slack.com/api/"
	SLACKPOSTMESSAGE = SLACKAPIURL + "chat.postMessage"
	SLACKPOSTFILE    = SLACKAPIURL + "files.upload"
	TOKEN            = ""
)

func (t *SlackApi) apipost() {

}

func (t *SlackApi) postSlackFileData(filename string, file *os.File) string {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return ""
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return ""
	}
	err = writer.Close()
	if err != nil {
		return ""
	}

	values := url.Values{
		"token": {t.Token},
	}
	values.Add("channels", t.Channel)
	values.Add("filename", filename)

	req, err := http.NewRequest("POST", SLACKPOSTFILE, buf)
	if err != nil {
		return ""
	}
	req = req.WithContext(context.Background())
	req.URL.RawQuery = (values).Encode()
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	// client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body2)
}

func (t *SlackApi) postSlackfile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return ""
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return ""
	}
	err = writer.Close()
	if err != nil {
		return ""
	}

	values := url.Values{
		"token": {t.Token},
	}
	values.Add("channels", t.Channel)
	values.Add("filename", filepath.Base(filename))

	req, err := http.NewRequest("POST", SLACKPOSTFILE, body)
	if err != nil {
		return ""
	}
	req = req.WithContext(context.Background())
	req.URL.RawQuery = (values).Encode()
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	// client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body2)
}

func (t *SlackApi) postSlackMessage(text string) string {
	urldata := SLACKPOSTMESSAGE
	values := url.Values{}
	values.Set("token", t.Token)
	values.Add("channel", t.Channel)
	values.Add("text", text)

	client := &http.Client{}
	client.Timeout = time.Second * 15
	req, err := http.NewRequest("POST", urldata, strings.NewReader(values.Encode()))
	if err != nil {
		return ""
	}
	//
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)

}
