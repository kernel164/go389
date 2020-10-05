package backend

import (
	"errors"

	"github.com/kernel164/go389/internal/backend/proxy"
	"github.com/kernel164/go389/internal/backend/yaml"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"
)

func GetBackendHandler(name string, args *model.ServerArgs) (model.BackendHandler, error) {
	switch cfg.GetBackendType(name) {
	case "yaml":
		return yaml.NewYamlBackendHandler(name, args)
	case "proxy":
		return proxy.NewProxyBackendHandler(name)
	}
	return nil, errors.New("backend not supported")
}
