package main

import (
	"context"
	"github.com/sebasmannem/tplr/internal"
)

var (
	ctx context.Context
	//ctxCancelFunc context.CancelFunc
	config internal.Config
)

func initContext() {
	ctx, _ = config.GetTimeoutContext(context.Background())
	internal.InitContext(ctx)
}

func main() {
	var err error
	initLogger()
	if config, err = internal.NewConfig(); err != nil {
		log.Fatal(err)
	} else {
		initRemoteLoggers()
		enableDebug(config.Debug)
		log.Debug("initializing config object")
		defer log.Sync() //nolint:errcheck
		log.Debug("checking if patching is required")
		initContext()
		if err = config.Tekton.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
