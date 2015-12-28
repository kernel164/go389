// +build !yaml

package yaml

import (
	model "../../model"
	"errors"
)

func NewYamlBackendHandler(name string) (model.BackendHandler, error) {
	return nil, errors.New("yaml backend not supported. Try build with yaml tag")
}
