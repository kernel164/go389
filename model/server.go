package model

type BaseServer struct {
	Base
}

type ServerHandler interface {
	Start(async bool) error
}
