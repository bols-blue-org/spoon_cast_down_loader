package main

import (
	"log"
	"os"

	"github.com/bols-blue-org/spoon_cast_downloader/src/spoon/cast"
	"github.com/urfave/cli"
)

var (
	argFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "download cast id.",
			Value: "",
		},
	}
)

func main() {
	app := cli.NewApp()

	app.Flags = argFlags
	app.Action = func(c *cli.Context) error {
		cast.Download(c)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	os.Exit(0)

}
