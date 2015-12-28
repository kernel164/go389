package server

import (
	cfg "../cfg"
	model "../model"
	ldap "./ldap"
	"errors"
)

func GetServerHandler(name string, backendHandler model.BackendHandler) (model.ServerHandler, error) {
	switch cfg.GetServerType(name) {
	case "ldap":
		return ldap.NewLdapServerHandler(name, backendHandler)
	}
	return nil, errors.New("server not supported")
}
