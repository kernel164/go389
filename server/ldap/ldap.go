package ldap

import (
	cfg "../../cfg"
	log "../../log"
	model "../../model"
	"github.com/nmcclain/ldap"
)

type LdapServerHandler struct {
	model.ServerHandler
	cfg     LdapServerSettings
	backend model.BackendHandler
}

type LdapServerSettings struct {
	model.BaseServer
	Bind     string
	CertFile string
	KeyFile  string
}

func NewLdapServerHandler(name string, backendHandler model.BackendHandler) (model.ServerHandler, error) {
	settings := LdapServerSettings{}
	cfg.GetServerCfg(name, &settings)
	return LdapServerHandler{cfg: settings, backend: backendHandler}, nil
}

func (h LdapServerHandler) Start(async bool) error {
	server := ldap.NewServer()
	server.EnforceLDAP = true
	server.BindFunc("", h.backend)
	server.SearchFunc("", h.backend)
	server.CloseFunc("", h.backend)

	switch h.cfg.Type {
	case "ldaps":
		log.Info("LDAPS server", "bind", h.cfg.Bind)
		if err := server.ListenAndServeTLS(h.cfg.Bind, h.cfg.CertFile, h.cfg.KeyFile); err != nil {
			log.Error("LDAP Server Failed: %s", err.Error())
		}
	default:
		log.Info("LDAP server", "bind", h.cfg.Bind)
		if err := server.ListenAndServe(h.cfg.Bind); err != nil {
			log.Error("LDAP Server Failed: %s", err.Error())
		}
	}

	return nil
}
