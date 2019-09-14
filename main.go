package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/urfave/cli"
)

const (
	volumeIdentifier = "cidata"
	fileMode         = 0700
	blockSize        = int64(2048)
)

func main() {
	app := cli.NewApp()
	app.Name = "quiso"
	app.Usage = "Quick Cloudinit ISOs"
	app.Action = build
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user-data, u",
			Value: "user-data.yaml",
			Usage: "user-data file to include",
		},
		cli.StringFlag{
			Name:  "meta-data, m",
			Value: "meta-data.yaml",
			Usage: "meta-data file to include",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "user-data.iso",
			Usage: "output iso file",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func build(c *cli.Context) error {
	userdata := c.String("user-data")
	metadata := c.String("meta-data")
	output := c.String("output")
	perm := os.FileMode(fileMode)

	diskImg := filepath.Join(output)
	iso, err := os.OpenFile(diskImg, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer iso.Close()

	fs, err := iso9660.Create(iso, 0, 0, blockSize)
	fs.Mkdir("/")
	if err != nil {
		return err
	}

	for filename, filepath := range map[string]string{"user-data": userdata, "meta-data": metadata} {
		f, err := ioutil.ReadFile(filepath) // just pass the file name
		if err != nil {
			return err
		}

		rw, err := fs.OpenFile("/"+filename, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return err
		}

		rw.Write(f)
	}

	if err = fs.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		VolumeIdentifier: volumeIdentifier,
	}); err != nil {
		return err
	}

	return nil
}
