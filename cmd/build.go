package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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

func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "",
		Action:  Build,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user-data",
				Aliases: []string{"u"},
				Value:   "",
				Usage:   "File containing User Data",
			},
			&cli.StringFlag{
				Name:    "meta-data",
				Aliases: []string{"m"},
				Value:   "",
				Usage:   "File containing Meta Data",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   defaultOutput,
				Usage:   "File to output final ISO file",
			},
		},
	}
}

func Build(c *cli.Context) error {
	userdata, metadata, output := getInputs(c)

	if userdata == "" {
		fmt.Errorf("no user data input, please specify a file with --user-data")
	}

	if userdata == "" {
		fmt.Errorf("no user data input, please specify a file with --user-data")
	}

	return buildIso(userdata, metadata, output)
}

func getInputs(c *cli.Context) (userdata, metadata, output string) {
	userdata = c.String("user-data")
	if userdata == "" {
		logrus.Debug("No --user-data param, checking for defaults")
		for _, ud := range possibleUserDataFiles {
			logrus.Debugf("Checking for %s", ud)
			if fileExists(ud) {
				logrus.Debugf("Found user data file: %s", ud)
				userdata = ud
				break
			}
		}
	}
	logrus.Debugf("User data input: %s", userdata)
	if !fileExists(userdata) {
		logrus.Debugf("User data file does not exist: %s", userdata)
		userdata = ""
	}

	metadata = c.String("meta-data")
	if metadata == "" {
		logrus.Debug("No --meta-data param, checking for defaults")
		for _, md := range possibleMetaDataFiles {
			logrus.Debugf("Checking for %s", md)
			if fileExists(md) {
				logrus.Debugf("Found user data file: %s", md)
				metadata = md
				break
			}
		}
	}
	logrus.Debugf("Meta data input: %s", metadata)
	if !fileExists(metadata) {
		logrus.Debugf("Meta data file does not exist: %s", metadata)
		metadata = ""
	}

	output = c.String("output")
	if output == "" {
		// fail safe if run as default action
		output = defaultOutput
	}
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

func buildIso(userdata, metadata, output string) error {
	perm := os.FileMode(fileMode)
	iso, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer iso.Close()

	fs, err := iso9660.Create(iso, 0, 0, blockSize, "")
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

	logrus.Infof("Cloud init ISO created: %s", output)
	return nil
}
