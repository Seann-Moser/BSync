package main

import (
	"github.com/Seann-Moser/BaseGoAPI/internal/server"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var (
		err error
		app = &cli.App{
			Name:   "",
			Usage:  "",
			Action: server.Serve,
		}
	)

	if err = app.Run(os.Args); err != nil {
		log.Println("exiting with error:", err.Error())
	}
}
