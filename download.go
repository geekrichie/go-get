package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var errNotValidType = errors.New("not valid file type")
var errDataIncomplete = errors.New("data is not complete")

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

func (d *Download)DownloadTrunk(start, end int) error{
	req, err := NewRequest("GET", d.url)
	if err != nil {
		return err
	}
	req.Header.Set("Range", "bytes="+strconv.Itoa(start) +"-"+strconv.Itoa(end))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(data) != end-start+1 {
		return errDataIncomplete
	}
	filename := GetFileNameFromResponce(resp)
	if filename == "" {
		return errNotValidType
	}
	f := &File{
		Name: filename,
		Size: len(data),
	}
	err = f.SaveChunk(int64(start),data)
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

var (
	MIN_CHUNK_SIZE  = 4
	MAX_CHUNK_SIZE  = 20*1024*1024
)


func (d *Download)GetChunkSize() int{
	if d.size == 0 {
		return 0
	}
	var chunkSize = d.size/runtime.NumCPU()
	for chunkSize > MAX_CHUNK_SIZE {
		chunkSize = chunkSize / 2
	}
	if chunkSize < MIN_CHUNK_SIZE {
		chunkSize = MIN_CHUNK_SIZE
	}
	return chunkSize
}


func (d *Download)DownloadMulti() {
	var (
		chunkSize int
		n int
		wg sync.WaitGroup
	)
	chunkSize = d.GetChunkSize()
	if d.size <=0 {
		return
	}
	n = d.size/chunkSize
	if d.size %chunkSize != 0 {
		n = n+1
	}

	for i:=0;i < n;i++ {
		wg.Add(1)
          if i!= n-1 {
          	go func(i int) {
				err := d.DownloadTrunk(i*chunkSize, (i+1)*chunkSize-1)
				if err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}(i)
		  }else {
		  	go func(i int) {
				err := d.DownloadTrunk(i*chunkSize, d.size-1)
				if err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}(i)
		  }
	}
	wg.Wait()
}