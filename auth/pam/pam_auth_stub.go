// +build !pam

package pam

import (
	"errors"
)

func PAMAuth(serviceName string, userName string, passwd string) error {
	return errors.New("PAM not supported. Try build with pam tag")
}
