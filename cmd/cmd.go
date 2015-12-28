package cmd

import (
	"gopkg.in/codegangsta/cli.v1"
)

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
