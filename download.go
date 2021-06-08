package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var errNotValidType = errors.New("not valid file type")

type Download struct {
	url string
	rangeable bool
	size int //下载文件的大小
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

func (d *Download) GetRangeInfo() error{
	req, err := NewRequest("GET", d.url)
	if err != nil {
		return err
	}
	req.Header.Set("Range", "bytes=0-1")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//PrintMap(resp.Header)
	if resp.Header.Get("Accept-Ranges") == "bytes"{
		d.rangeable = true
	}
	if resp.Header.Get("Content-Range") != "" {
		d.rangeable = true
		contentRange := resp.Header.Get("Content-Range")
		pos := strings.LastIndexByte(contentRange, '/')
		size := contentRange[pos+1:]
		size = strings.TrimSpace(size)
		if size != "*"{
			d.size, err = strconv.Atoi(size)
		}
	}
	return nil
}
