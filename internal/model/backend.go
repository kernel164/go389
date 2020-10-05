package model

import (
	"github.com/nmcclain/ldap"
)

type BaseBackend struct {
	Base
	BaseDN string
}

type BackendHandler interface {
	ldap.Binder
	ldap.Searcher
	ldap.Closer
}
