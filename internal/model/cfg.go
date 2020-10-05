package model

type Base struct {
	Type string
}

type BaseCfg struct {
	Base
}

const ProgramName = "go389"

type CfgHandler interface {
	GetServer() string
	GetBackend() string
	GetServerType(name string) string
	GetServerCfg(name string, x interface{})
	GetBackendType(name string) string
	GetBackendCfg(name string, x interface{})
	GetAuthType(name string) string
	GetAuthCfg(name string, x interface{})
}

type HashArgs struct {
	Algo  string `arg:"-a,--algo" default:"sha256"`
	Value string `arg:"-v,--value"`
}

type ServerArgs struct {
	Config string `arg:"-c,--config,required"`
}
