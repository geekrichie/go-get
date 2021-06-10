package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (f *File)SaveChunk(data []byte) error{
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


func PrintMap(data interface{}) {
	jsonSerilize, _:= json.Marshal(data)
	fmt.Println(string(jsonSerilize))
}
