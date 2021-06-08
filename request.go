package main

import (
	"go-get/mimetype"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)


func NewRequest(method string,url string) (*http.Request,error){
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return req,err
	}
	req.Header.Set("User-Agent","go-get/1.0")
	return req, err
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
		typeMap := mimetype.GetMimeTypeMap()
		if _,ok := typeMap[mediaType] ; ok {
			filename = "index." + typeMap[mediaType]
		}
	}else {
		filename = url[i+1:]
	}
	return filename
}