package main

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

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
		return errors.New("")
	}else if err != nil {
		return err
	}
	 return nil
}


func callerFileLine() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	return
}