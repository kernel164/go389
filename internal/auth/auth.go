package auth

import (
	"errors"

	"github.com/kernel164/go389/internal/auth/pam"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
	"github.com/kernel164/go389/internal/util"
)

func GetAuthHandler(name string) (model.AuthHandler, error) {
	switch cfg.GetAuthType(name) {
	case "pam":
		return pam.NewPamAuthHandler(name)
	}
	return NewDefaultAuthHandler(name)
}

type DefaultAuthHandler struct {
	model.AuthHandler
	settings DefaultAuthSettings
}

type DefaultAuthSettings struct {
	model.BaseAuth
	Hash string
}

func NewDefaultAuthHandler(name string) (model.AuthHandler, error) {
	settings := DefaultAuthSettings{}
	cfg.GetAuthCfg(name, &settings)
	return DefaultAuthHandler{settings: settings}, nil
}

func (h DefaultAuthHandler) Auth(userName string, backendPasswd string, checkPasswd string) error {
	switch h.settings.Hash {
	case "md5":
		ok, err := util.CheckMd5(checkPasswd, backendPasswd)
		if err != nil || !ok {
			return errors.New("invalid password")
		}
	case "bcrypt":
		ok, err := util.CheckBcrypt(checkPasswd, backendPasswd)
		if err != nil || !ok {
			return errors.New("invalid password")
		}
	case "sha256":
		ok, err := util.CheckSha256(checkPasswd, backendPasswd)
		if err != nil || !ok {
			return errors.New("invalid password")
		}
	}
	return nil
}
