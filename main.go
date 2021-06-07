package main

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"log"
	"net/http"
	"os"
	"strings"
)

var url string

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Usage:   "a web url",
				Value: "http://localhost",
				Aliases: []string{"b"},
				Destination: &url,
			},
		},
		Action: func(context *cli.Context) error {
			return DownloadSingle()
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadSingle() (err error){
	var data []byte
	var resp *http.Response

	if !strings.HasPrefix(url,"http://") || !strings.HasPrefix(url,"https://") {
		url  = "http://" + url
	}

	resp,data, err = NewRequestWithResponse(url)
	if err != nil {
		return err
	}
	f := &File{
		Name: GetFileNameFromResponce(resp),
		Size: len(data),
	}
	err = f.Save(data)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}