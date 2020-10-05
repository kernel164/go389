package cmd

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/model"
	"github.com/kernel164/go389/internal/util"
)

func RunHash(log logr.Logger, args *model.HashArgs) error {
	value := args.Value
	if value == "" {
		return nil
	}
	switch args.Algo {
	case "sha256":
		fmt.Println(util.Sha256(value))
	case "bcrypt":
		v, err := util.Bcrypt(value)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		fmt.Println(v)
	case "md5":
		fmt.Println(util.Md5(value))
	}

	return nil
}
