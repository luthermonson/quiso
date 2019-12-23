package cmd

import "github.com/urfave/cli"

func Commands() cli.Commands {
	return []cli.Command{
		BuildCommand(),
	}
}
