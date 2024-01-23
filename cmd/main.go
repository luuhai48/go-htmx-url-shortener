package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cli := &cli.App{
		Name:                 "URL Shortener",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Action:               startServer,
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
