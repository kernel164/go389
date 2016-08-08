package cmd

import (
	"path/filepath"

	backend "../backend"
	cfg "../cfg"
	log "../log"
	server "../server"
	"gopkg.in/codegangsta/cli.v1"
)

var CmdServer = cli.Command{
	Name:        "server",
	Usage:       "Server",
	Description: `Server`,
	Action:      runServer,
	Flags: []cli.Flag{
		stringFlag("config, c", "app.yml", "configuration file path"),
	},
}

func runServer(c *cli.Context) error {
	config := c.String("config")
	extn := filepath.Ext(config)
	extn = extn[1:]

	if err := log.Init(""); err != nil {
		return err
	}

	//log.Info("Loading config", "file", config, "type", extn)

	// load config
	if err := cfg.Load(extn, config); err != nil {
		log.Error(err.Error())
		return err
	}

	// backend
	backendHandler, err := backend.GetBackendHandler(cfg.GetBackend())
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// server
	serverHandler, err := server.GetServerHandler(cfg.GetServer(), backendHandler)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// start server
	serverHandler.Start(false)
	return nil
}
