package server

import (
	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
	"github.com/nmcclain/ldap"
)

type handler struct {
	model.ServerHandler
	log     logr.Logger
	cfg     *Config
	backend model.BackendHandler
}

// Config - server settings
type Config struct {
	model.BaseServer
	Bind        string
	EnforceLDAP bool
	CertFile    string
	KeyFile     string
}

// New - create new server
func New(name string, log logr.Logger, backendHandler model.BackendHandler) (model.ServerHandler, error) {
	config := &Config{}
	cfg.GetServerCfg(name, config)
	return &handler{
		log:     log,
		cfg:     config,
		backend: backendHandler,
	}, nil
}

func (h *handler) Start(async bool) error {
	server := ldap.NewServer()
	server.EnforceLDAP = h.cfg.EnforceLDAP
	server.BindFunc("", h.backend)
	server.SearchFunc("", h.backend)
	server.CloseFunc("", h.backend)

	switch h.cfg.Type {
	case "ldaps":
		h.log.Info("LDAPS server", "bind", h.cfg.Bind)
		if err := server.ListenAndServeTLS(h.cfg.Bind, h.cfg.CertFile, h.cfg.KeyFile); err != nil {
			h.log.Info("LDAP Server Failed.", "err", err)
			return err
		}
	default:
		h.log.Info("LDAP server", "bind", h.cfg.Bind)
		if err := server.ListenAndServe(h.cfg.Bind); err != nil {
			h.log.Info("LDAP Server Failed.", "err", err)
			return err
		}
	}

	return nil
}
