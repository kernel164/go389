package yaml

import (
	"io/ioutil"
	"reflect"

	"github.com/kernel164/go389/internal/model"

	"gopkg.in/yaml.v3"
)

type handler struct {
	model.CfgHandler
	config *Config
}

// Config - config
type Config struct {
	Server   string
	Backend  string
	Backends map[string]map[string]interface{}
	Servers  map[string]map[string]interface{}
	Auths    map[string]map[string]interface{}
}

// New - config handler
func New(file string) (model.CfgHandler, error) {
	yamlCfg := &Config{}
	content, err := ioutil.ReadFile(file)
	if err == nil {
		if unmarshalerr := yaml.Unmarshal(content, yamlCfg); unmarshalerr != nil {
			panic(unmarshalerr)
			return nil, unmarshalerr
		}
	} else {
		return nil, err
	}
	return handler{config: yamlCfg}, nil
}

func (h handler) GetServer() string {
	return h.config.Server
}

func (h handler) GetBackend() string {
	return h.config.Backend
}

func (h handler) GetServerType(name string) string {
	if val, ok := h.config.Servers[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h handler) GetServerCfg(name string, x interface{}) {
	mapToStruct(h.config.Servers[name], x)
}

func (h handler) GetBackendType(name string) string {
	if val, ok := h.config.Backends[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h handler) GetBackendCfg(name string, x interface{}) {
	mapToStruct(h.config.Backends[name], x)
}

func (h handler) GetAuthType(name string) string {
	if val, ok := h.config.Auths[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h handler) GetAuthCfg(name string, x interface{}) {
	mapToStruct(h.config.Auths[name], x)
}

func mapToStruct(m map[string]interface{}, x interface{}) {
	obj := reflect.ValueOf(x).Elem()
	for key, value := range m {
		field := obj.FieldByName(key)
		if !field.IsValid() {
			// or handle as error if you don't expect unknown values
			continue
		}
		if !field.CanSet() {
			// or return an error on private fields
			continue
		}
		field.Set(reflect.ValueOf(value))
	}
}
