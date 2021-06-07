package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 20,
		IdleConnTimeout: 1000,
	},
}

func NewRequest(url string) ([]byte,error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil,err
	}
	resp, err := client.Do(req)
	if err != nil{
		return nil,err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil,err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	return content,nil
}

func NewRequestWithResponse(url string) (*http.Response,[]byte, error ){
	var resp *http.Response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp,nil,err
	}
	resp, err = client.Do(req)
	if err != nil{
		return resp,nil,err
	}
	defer resp.Body.Close()
	if err != nil {
		return resp,nil,err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp,nil,err
	}
	return resp,content,nil
}


func GetFileNameFromResponce(resp *http.Response) string{
	content := resp.Header.Get("Content-Disposition")
	_,params, _ := mime.ParseMediaType(content)
	if _, ok := params["filename"]; ok {
		return params["filename"]
	}
	url := resp.Request.URL.Path
	i := strings.LastIndexByte(url, '/')
	var filename string
	if i == len(url)-1 {
		content = resp.Header.Get("Content-Type")
		mediaType,_, _ := mime.ParseMediaType(content)
		if mediaType == "text/html" {
			filename = "index.html"
		}
	}else {
		filename = url[i+1:]
	}
	fmt.Println(filename)
	return filename
}