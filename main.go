package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/luthermonson/quiso/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "quiso"
	app.Usage = "Quick CloudInit ISOs"
	app.Commands = cmd.Commands()
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if ctx.GlobalBool("quiet") {
			logrus.SetOutput(ioutil.Discard)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Turn on verbose debug logging",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Turn on off all logging",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
