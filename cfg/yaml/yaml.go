// +build yaml

package yaml

import (
	model "../../model"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"reflect"
)

type YamlCfgHandler struct {
	model.CfgHandler
	file YamlCfg
}

type YamlCfg struct {
	Server   string
	Backend  string
	Backends map[string]map[string]interface{}
	Servers  map[string]map[string]interface{}
	Auths    map[string]map[string]interface{}
}

func NewYamlCfgHandler(name string, file string) (model.CfgHandler, error) {
	yamlCfg := YamlCfg{}
	content, err := ioutil.ReadFile(file)
	if err == nil {
		if unmarshalerr := yaml.Unmarshal(content, &yamlCfg); unmarshalerr != nil {
			panic(unmarshalerr)
			return nil, unmarshalerr
		}
	} else {
		return nil, err
	}
	return YamlCfgHandler{file: yamlCfg}, nil
}

func (h YamlCfgHandler) GetServer() string {
	return h.file.Server
}

func (h YamlCfgHandler) GetBackend() string {
	return h.file.Backend
}

func (h YamlCfgHandler) GetServerType(name string) string {
	if val, ok := h.file.Servers[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h YamlCfgHandler) GetServerCfg(name string, x interface{}) {
	mapToStruct(h.file.Servers[name], x)
}

func (h YamlCfgHandler) GetBackendType(name string) string {
	if val, ok := h.file.Backends[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h YamlCfgHandler) GetBackendCfg(name string, x interface{}) {
	mapToStruct(h.file.Backends[name], x)
}

func (h YamlCfgHandler) GetAuthType(name string) string {
	if val, ok := h.file.Auths[name]["Type"].(string); ok {
		return val
	}
	return name
}

func (h YamlCfgHandler) GetAuthCfg(name string, x interface{}) {
	mapToStruct(h.file.Auths[name], x)
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
