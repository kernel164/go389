package cfg

import (
	"github.com/kernel164/go389/internal/cfg/yaml"
	"github.com/kernel164/go389/internal/model"
)

var h model.CfgHandler

func Load(file string) error {
	cfgHandler, err := yaml.New(file)
	if err != nil {
		return err
	}
	h = cfgHandler
	return nil
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
