package adapter

import (
	"github.com/n-creativesystem/oidc-proxy/adapter/command"
	"github.com/n-creativesystem/oidc-proxy/version"

	"github.com/urfave/cli/v2"
)

type Adapter interface {
	Run(args []string) error
}

func New() Adapter {
	app := cli.NewApp()
	app.Name = "openid connect proxy server"
	app.Version = version.Version
	app.Description = "openid connect proxy server"
	app.Commands = []*cli.Command{
		command.ProxyCommand,
		command.AppFileCommand,
	}
	return app
}
