// +build !ini

package ini

import (
	model "../../model"
	"errors"
)

func NewIniCfgHandler(name string, file string) (model.CfgHandler, error) {
	return nil, errors.New("ini cfg not supported. Try build with ini tag")
}
