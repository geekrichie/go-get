package main

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var errNotValidType = errors.New("not valid file type")

type Download struct {
	url string
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 20,
		IdleConnTimeout: 1000,
	},
}

func (d *Download)downloadFull() error{
	req, err := NewRequest("GET", d.url)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	filename := GetFileNameFromResponce(resp)
	if filename == "" {
		return errNotValidType
	}
	f := &File{
		Name: filename,
		Size: len(data),
	}
	err = f.Save(data)
	return err
}
