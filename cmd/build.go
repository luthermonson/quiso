package cmd

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/urfave/cli"
)

const (
	volumeIdentifier = "cidata"
	fileMode         = 0700
	blockSize        = int64(2048)
	defaultOutput    = "user-data.iso"
)

// if user-data not passed check these defaults
var possibleUserDataFiles = []string{
	"user-data.yaml",
	"user-data.yml",
	"user-data",
}

var possibleMetaDataFiles = []string{
	"meta-data.yaml",
	"meta-data.yml",
	"meta-data",
}

func BuildCommand() cli.Command {
	return cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "",
		Action:  build,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "user-data, u",
				Value: "",
				Usage: "File containing User Data",
			},
			cli.StringFlag{
				Name:  "meta-data, m",
				Value: "",
				Usage: "File containing Meta Data",
			},
			cli.StringFlag{
				Name:  "output, o",
				Value: defaultOutput,
				Usage: "File to output final ISO file",
			},
		},
	}
}

func build(c *cli.Context) error {
	userdata, metadata, output := getInputs(c)

	perm := os.FileMode(fileMode)
	iso, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, perm)
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
		f, err := ioutil.ReadFile(filepath)
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

func getInputs(c *cli.Context) (userdata, metadata, output string) {
	userdata = c.String("user-data")
	if userdata == "" {
		for _, ud := range possibleUserDataFiles {
			if fileExists(ud) {
				userdata = ud
				break
			}
		}
	}
	logrus.Debugf("User data input: %s", userdata)

	metadata = c.String("meta-data")
	if metadata == "" {
		for _, md := range possibleMetaDataFiles {
			if fileExists(md) {
				metadata = md
				break
			}
		}
	}
	logrus.Debugf("Meta data input: %s", metadata)

	output = c.String("output")
	logrus.Debugf("Output ISO: %s", output)

	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
