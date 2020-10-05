// +build ini

package ini

import (
	"github.com/kernel164/go389/internal/model"

	"gopkg.in/ini.v1"
)

type IniCfgHandler struct {
	model.CfgHandler
	file *ini.File
}

func NewIniCfgHandler(name string, file string) (model.CfgHandler, error) {
	iniFile, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	iniFile.NameMapper = ini.AllCapsUnderscore
	return IniCfgHandler{file: iniFile}, nil
}

func (h IniCfgHandler) GetServer() string {
	return h.file.Section("").Key("Server").MustString("ldap")
}

func (h IniCfgHandler) GetBackend() string {
	return h.file.Section("").Key("Backend").MustString("yaml")
}

func (h IniCfgHandler) GetServerType(name string) string {
	return h.file.Section("server." + name).Key("Type").MustString(name)
}

func (h IniCfgHandler) GetServerCfg(name string, x interface{}) {
	h.file.Section("server." + name).MapTo(x)
}

func (h IniCfgHandler) GetBackendType(name string) string {
	return h.file.Section("backend." + name).Key("Type").MustString(name)
}

func (h IniCfgHandler) GetBackendCfg(name string, x interface{}) {
	h.file.Section("backend." + name).MapTo(x)
}

func (h IniCfgHandler) GetAuthType(name string) string {
	return h.file.Section("auth." + name).Key("Type").MustString(name)
}

func (h IniCfgHandler) GetAuthCfg(name string, x interface{}) {
	h.file.Section("auth." + name).MapTo(x)
}
