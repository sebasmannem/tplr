package internal

import (
	"context"
	"github.com/sebasmannem/tplr/pkg/tekton_handler"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
	//ctx context.Context
)

func InitLogger(logger *zap.SugaredLogger) {
	log = logger
	tekton_handler.InitLogger(log)
}

func InitContext(c context.Context) {
	tekton_handler.InitContext(c)
}
