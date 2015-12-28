package cfg

import (
	model "../model"
	ini "./ini"
	yaml "./yaml"
	"errors"
)

var h model.CfgHandler

func Load(name string, file string) error {
	cfgHandler, err := getCfgHandler(name, file)
	if err != nil {
		return err
	}
	h = cfgHandler
	return nil
}

func getCfgHandler(name string, file string) (model.CfgHandler, error) {
	switch name {
	case "ini":
		return ini.NewIniCfgHandler(name, file)
	case "yaml":
	case "yml":
		return yaml.NewYamlCfgHandler(name, file)
	}
	return nil, errors.New("cfg not supported")
}

func GetServer() string {
	return h.GetServer()
}

func GetBackend() string {
	return h.GetBackend()
}

func GetServerType(name string) string {
	return h.GetServerType(name)
}

func GetServerCfg(name string, x interface{}) {
	h.GetServerCfg(name, x)
}

func GetBackendType(name string) string {
	return h.GetBackendType(name)
}

func GetBackendCfg(name string, x interface{}) {
	h.GetBackendCfg(name, x)
}

func GetAuthType(name string) string {
	return h.GetAuthType(name)
}

func GetAuthCfg(name string, x interface{}) {
	h.GetAuthCfg(name, x)
}
