package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var  errLengthInvalid = errors.New("写入长度错误")

type File struct {
	Name string
	Size int
	Path string
}


func (f *File)Save(data []byte) error{
     fl,err := os.OpenFile(f.Name,os.O_CREATE|os.O_RDWR,0666)
     defer fl.Close()
     if err != nil {
     	return err
	 }
	n, err := fl.Write(data)
	if n != len(data) {
		return errLengthInvalid
	}else if err != nil {
		return err
	}
	 return nil
}

func (f *File)SaveChunk(r io.Reader) error{
	if !Exists(f.Path) {
		err := os.MkdirAll(f.Path, 0664)
		CheckError(err)
	}
	fl,err := os.OpenFile(f.Path+string(os.PathSeparator)+f.Name,os.O_CREATE|os.O_RDWR,0666)
	defer fl.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(fl, r)
	if err != nil{
		return err
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func PrintMap(data interface{}) {
	jsonSerilize, _:= json.Marshal(data)
	fmt.Println(string(jsonSerilize))
}
