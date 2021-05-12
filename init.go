package main

import (
	"log"
	"os"

	"github.com/n-creativesystem/oidc-proxy/logger"
)

func init() {
	// if err := envconfig.Process("", &env.Env); err != nil {
	// 	log.Fatalln(err)
	// }

	logger.Log = logger.New(os.Stdout, logger.Info, logger.FormatStandard, logger.FormatDate)
	log.SetOutput(logger.Log)
}
