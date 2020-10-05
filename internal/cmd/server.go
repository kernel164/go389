package cmd

import (
	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/backend"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
	"github.com/kernel164/go389/internal/server"
)

func RunServer(log logr.Logger, args *model.ServerArgs) error {
	// load config
	if err := cfg.Load(args.Config); err != nil {
		log.Error(err, "error loading config")
		return err
	}

	// backend
	backendHandler, err := backend.GetBackendHandler(cfg.GetBackend(), log, args)
	if err != nil {
		log.Error(err, "error getting backend")
		return err
	}

	// server
	serverHandler, err := server.New(cfg.GetServer(), log, backendHandler)
	if err != nil {
		log.Error(err, "error getting server")
		return err
	}

	// start server
	serverHandler.Start(false)
	return nil
}
