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
	rangeable bool //是否可以分块下载
	size int       //下载文件的大小
	chunkSize int  //分块下载文件的大小
	chunkNum   int //分块下载文件的数量
	finished   int //下载完成的分区数量
	retry      int //重试次数
	dir        string //文件下载目录
	filename   string //文件名称
}

type Chunk struct {
	close chan struct{} //是否已经完成
	size int            //当前chunk的大小
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 20,
		IdleConnTimeout: 1000,
	},
}

//设置chunkSize的大小
func (d *Download)SetChunkSize(chunkSize int) {
	d.chunkSize = chunkSize
}
//设置chunkSize的大小
func (d *Download)SetChunkNum(chunkNum int) {
	d.chunkNum = chunkNum
}

//下载整个文件
func (d *Download)DownloadFull() error{
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

//下载某个chunk
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
		Name: fmt.Sprintf("%s_%d", filename, start),
		Size: len(data),
		//Path: fmt.Sprintf("%s%s_%d", d.dir,filename,string(os.PathSeparator), start),
	}
	err = f.SaveChunk(data)
	return err
}

//func (d *Download) combineFilePath() string{
//	return fmt.Sprintf("%s%s_%d",d.dir)
//}

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
	MIN_CHUNK_SIZE  = 5*1024*1024
	MAX_CHUNK_SIZE  = 20*1024*1024
)


//计算下载的文件块大小
func (d *Download)GetChunkSize() int{
	if d.size == 0 {
		return 0
	}
	var chunkSize = d.size/runtime.NumCPU()
	for chunkSize > MAX_CHUNK_SIZE {
		chunkSize = chunkSize / 2
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
	d.SetChunkSize(chunkSize)
	if d.size <=0 {
		return
	}
	n = d.size/chunkSize
	if d.size %chunkSize != 0 {
		n = n+1
	}
	d.SetChunkNum(n)

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