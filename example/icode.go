package main

import (
	"os"

	"github.com/it-chain/sdk"
	"github.com/it-chain/sdk/example/handler"
	"github.com/it-chain/sdk/logger"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `short:"p" long:"port" description:"set port"`
}

func main() {
	logger.EnableFileLogger(true, "./icode.log")
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		logger.Error(nil, "fail parse args: "+err.Error())
		os.Exit(1)
	}

	exHandler := &handler.HandlerExample{}
	ibox := sdk.NewIBox(opts.Port)
	ibox.SetHandler(exHandler)
	ibox.On(30)
}
