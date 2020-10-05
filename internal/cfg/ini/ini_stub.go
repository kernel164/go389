// +build !ini

package ini

import (
	"errors"

	"github.com/kernel164/go389/internal/model"
)

func NewIniCfgHandler(name string, file string) (model.CfgHandler, error) {
	return nil, errors.New("ini cfg not supported. Try build with ini tag")
}
