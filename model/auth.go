package model

type BaseAuth struct {
	Base
}

type AuthHandler interface {
	Auth(userName string, backendPasswd string, checkPasswd string) error
}
