package main

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"log"
	"os"
	"strings"
)

var url string
var dir string

func main() {
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
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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