package server

import (
	"errors"

	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
	"github.com/kernel164/go389/internal/server/ldap"
)

func GetServerHandler(name string, backendHandler model.BackendHandler) (model.ServerHandler, error) {
	switch cfg.GetServerType(name) {
	case "ldap":
		return ldap.NewLdapServerHandler(name, backendHandler)
	}
	return nil, errors.New("server not supported")
}
