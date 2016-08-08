package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"gopkg.in/codegangsta/cli.v1"
)

var CmdHash = cli.Command{
	Name:        "hash",
	Usage:       "hash",
	Description: `hash`,
	Action:      runHash,
	Flags: []cli.Flag{
		stringFlag("algo, a", "sha256", "algo"),
		stringFlag("value, v", "", ""),
	},
}

func runHash(c *cli.Context) error {
	value := c.String("value")
	if value == "" {
		return nil
	}
	switch c.String("algo") {
	case "sha256":
		hash := sha256.New()
		hash.Write([]byte(value))
		fmt.Println(hex.EncodeToString(hash.Sum(nil)))
	}

	return nil
}
