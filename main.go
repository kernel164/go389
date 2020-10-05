package main

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/kernel164/go389/internal/cmd"
	"github.com/kernel164/go389/internal/model"
)

type args struct {
	Hash   *model.HashArgs   `arg:"subcommand:hash"`
	Server *model.ServerArgs `arg:"subcommand:server"`
}

func main() {
	args := &args{}
	p := arg.MustParse(args)
	switch {
	case args.Hash != nil:
		cmd.RunHash(args.Hash)
	case args.Server != nil:
		cmd.RunServer(args.Server)
	default:
		p.WriteHelp(os.Stdout)
	}
}
