package main

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"strings"
)

var url string
var dir string

func main() {
	var f *os.File
	var err error

		f, err = os.Create("cpu1.profile")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

	// ... rest of the program ...

		f, err = os.Create("mem1.profile")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Usage:   "a web url",
				Value: "localhost",
				Aliases: []string{"b"},
				Destination: &url,
			},
			&cli.StringFlag{
				Name:    "dir",
				Usage:   "file saved folder name",
				Value:    ".",
				Aliases: []string{"d"},
				Destination: &dir,
			},
		},
		Action: func(context *cli.Context) error {
			return prepareAction()
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func prepareAction() (err error){

	if !strings.HasPrefix(url,"http://") && !strings.HasPrefix(url,"https://") {
		url  = "http://" + url
	}
	d := &Download{
		url:url,
		done : make(chan struct{},1),
	}
	if dir != "" {
		d.dir = dir
	}
	err = d.GetRangeInfo()
	if err != nil {
		return err
	}
	//支持分段下载以及大小超过1M才分段下载
	if d.rangeable && d.size > MIN_CHUNK_SIZE{
		d.DownloadMulti()
	}else {
		d.DownloadFull()
	}


	return nil
}