package backend

import (
	"errors"

	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/backend/proxy"
	"github.com/kernel164/go389/internal/backend/yaml"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
)

func GetBackendHandler(name string, log logr.Logger, args *model.ServerArgs) (model.BackendHandler, error) {
	switch cfg.GetBackendType(name) {
	case "yaml":
		return yaml.New(name, log, args)
	case "proxy":
		return proxy.New(name, log)
	}
	return nil, errors.New("backend not supported")
}
