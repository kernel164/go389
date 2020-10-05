package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/go-logr/zapr"
	"github.com/kernel164/go389/internal/cmd"
	"github.com/kernel164/go389/internal/model"
	"go.uber.org/zap"
)

type args struct {
	Hash   *model.HashArgs   `arg:"subcommand:hash"`
	Server *model.ServerArgs `arg:"subcommand:server"`
}

func main() {
	// parse args
	args := &args{}
	p := arg.MustParse(args)

	// log
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log := zapr.NewLogger(zapLog)

	// run
	switch {
	case args.Hash != nil:
		cmd.RunHash(log, args.Hash)
	case args.Server != nil:
		cmd.RunServer(log, args.Server)
	default:
		p.WriteHelp(os.Stdout)
	}
}
