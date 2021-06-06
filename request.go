package main

import (
	"io/ioutil"
	"net/http"
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
