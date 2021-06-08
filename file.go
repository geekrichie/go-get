package main

import (
	"errors"
	"os"
)

var  errLengthInvalid = errors.New("写入长度错误")

type File struct {
	Name string
	Size int
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
