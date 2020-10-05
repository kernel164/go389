package cmd

import (
	"path/filepath"

	"github.com/kernel164/go389/internal/backend"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/log"
	"github.com/kernel164/go389/internal/model"
	"github.com/kernel164/go389/internal/server"
)

func RunServer(args *model.ServerArgs) error {
	config := args.Config
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
	backendHandler, err := backend.GetBackendHandler(cfg.GetBackend(), args)
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
