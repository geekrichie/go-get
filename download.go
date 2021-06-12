package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var errNotValidType = errors.New("not valid file type")
var errDataIncomplete = errors.New("data is not complete")

type Download struct {
	url string
	rangeable bool //是否可以分块下载
	size int       //下载文件的大小
	chunkSize int  //分块下载文件的大小
	chunkNum   int32 //分块下载文件的数量
	finished   int32 //下载完成的分区数量
	retry      int //重试次数
	dir        string //文件下载目录
	filename   string //文件名称
	done       chan struct{}
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
	d.chunkNum = int32(chunkNum)
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
	f := &File{
		Name: d.chunkName(start/d.chunkSize),
		Size: len(data),
		Path: d.combineFilePath(),
	}
	err = f.SaveChunk(data)
	if err == nil {
		//保存文件成功
         atomic.AddInt32(&d.finished, 1)
         if d.chunkNum == d.finished {
         	d.Close()
		 }
	}
	return err
}

var once sync.Once

func (d *Download) Close() {
	once.Do(func() {
		close(d.done)
	})
}

func (d *Download)chunkName(block int) string{
	return fmt.Sprintf("chunk_%d", block)
}

/**
合并后删除分区文件
 */
func (d *Download)clean() {
	path := d.combineFilePath()
	files,err := ioutil.ReadDir(path)
	CheckError(err)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "chunk_"){
			 os.Remove(path+string(os.PathSeparator)+file.Name())
		}
	}
}
/**
阻塞等待下载任务完成
合并文件
 */

func (d *Download)Merge() {
	ctx, cancel  := context.WithTimeout(context.Background(),time.Duration(10)*time.Minute)
	defer cancel()
	defer d.clean()
	select {
	    case <-ctx.Done():
	    	log.Println(d.filename,"下载已超时")
	    	log.Fatal("下载超时")
		case <-d.done:
			log.Println("准备合并操作")
	}
	path := d.combineFilePath()
	mergeFile,err := os.OpenFile(path+ string(os.PathSeparator) + d.filename, os.O_CREATE|os.O_WRONLY, 0664)
	CheckError(err)
	defer  mergeFile.Close()
	for  i  := int32(0); i< d.chunkNum; i++ {
		content, err := ioutil.ReadFile(path+ string(os.PathSeparator) + d.chunkName(int(i)))
		CheckError(err)
		n, err := mergeFile.Write(content)
		if n != len(content) {
			CheckError(errDataIncomplete)
		}
		CheckError(err)
	}
	log.Println("合并文件完成")
}

func (d *Download) combineFilePath() string{
	suffixPoint := strings.LastIndex(d.filename, ".")
	return fmt.Sprintf("%s%s%s",
		d.dir,
		string(os.PathSeparator),
		d.filename[:suffixPoint] )
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
	d.filename = GetFileNameFromResponce(resp)
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
		downloadFilePath string
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
	wg.Add(1)
	go func() {
		d.Merge()
		wg.Done()
	}()
	downloadFilePath = d.combineFilePath()

	for i:=0;i < n;i++ {
		//之前存在的文件就不重新下载了
		if !Exists(downloadFilePath + string(os.PathSeparator) + d.chunkName(i)) {
              continue
		}
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