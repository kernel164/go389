package backend

import (
	cfg "../cfg"
	model "../model"
	proxy "./proxy"
	yaml "./yaml"
	"errors"
)

func GetBackendHandler(name string) (model.BackendHandler, error) {
	switch cfg.GetBackendType(name) {
	case "yaml":
		return yaml.NewYamlBackendHandler(name)
	case "proxy":
		return proxy.NewProxyBackendHandler(name)
	}
	return nil, errors.New("backend not supported")
}
