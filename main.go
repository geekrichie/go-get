package main

import (
	"github.com/urfave/cli/v2" // imports as package "cli"
	"log"
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
				Value: "localhost",
				Aliases: []string{"b"},
				Destination: &url,
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

	if !strings.HasPrefix(url,"http://") || !strings.HasPrefix(url,"https://") {
		url  = "http://" + url
	}
	d := &Download{
		url:url,
	}
	d.downloadFull()


	return nil
}