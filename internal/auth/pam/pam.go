package pam

import (
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
)

type PamAuthHandler struct {
	model.AuthHandler
	settings PamAuthSettings
}

type PamAuthSettings struct {
	model.BaseAuth
	Service string
}

func NewPamAuthHandler(name string) (model.AuthHandler, error) {
	settings := PamAuthSettings{}
	cfg.GetAuthCfg(name, &settings)
	return PamAuthHandler{settings: settings}, nil
}

func (h PamAuthHandler) Auth(userName string, backendPasswd string, checkPasswd string) error {
	return PAMAuth(h.settings.Service, userName, checkPasswd)
}
