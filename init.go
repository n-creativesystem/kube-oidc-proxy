package main

import (
	"log"
	"os"

	"kube-oidc-proxy/env"
	"kube-oidc-proxy/logger"

	"github.com/kelseyhightower/envconfig"
)

func init() {
	if err := envconfig.Process("", &env.Env); err != nil {
		log.Fatalln(err)
	}

	logger.Log = logger.New(os.Stdout, env.Env.LogLevel)
}
