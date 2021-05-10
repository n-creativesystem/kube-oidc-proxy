package command

import (
	"crypto/tls"
	"net/http"

	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/watch"
	"github.com/n-creativesystem/oidc-proxy/watch/cert"
	watchConfig "github.com/n-creativesystem/oidc-proxy/watch/config"
	"github.com/urfave/cli/v2"
)

const (
	run     = "run"
	appConf = "config"
)

var ProxyCommand = &cli.Command{
	Name:  "proxy",
	Usage: "proxy server",
	Subcommands: []*cli.Command{
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Description: "proxy server run",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    appConf,
					Aliases: []string{"c"},
					Usage:   "application file path",
					Value:   "application.yaml",
				},
			},
			Action: func(c *cli.Context) error {
				https := true
				appWatcher, err := AppConfig(c.String(appConf))
				if err != nil {
					return err
				}
				appConf := appWatcher.Watching.(*watchConfig.Watch).Config
				cmWatcher, err := CertConfig(appConf.SslCertificate, appConf.SslCertificateKey)
				if err != nil {
					if err == watch.ErrFileNotFound {
						https = false
					} else {
						return err
					}
				}
				multiHost := appWatcher.Watching.(*watchConfig.Watch).MultiHost
				defer multiHost.Close()
				s := &http.Server{
					Addr:    appConf.GetPort(),
					Handler: multiHost,
				}
				logger.Log.Info("Application Start")
				if https {
					s.TLSConfig = &tls.Config{
						GetCertificate: cmWatcher.Watching.(*cert.Watch).GetCertificate,
					}
					return s.ListenAndServeTLS("", "")
				} else {
					return s.ListenAndServe()
				}
			},
		},
	},
}

func AppConfig(applicationFilePath string) (*watch.Watch, error) {
	appWatcher, err := watch.New(logger.Log)
	if err != nil {
		return nil, err
	}
	aw, err := watchConfig.New(applicationFilePath)
	if err != nil {
		return nil, err
	}
	appWatcher.Watching = aw
	if err := appWatcher.Watch(); err != nil {
		return nil, err
	}
	return appWatcher, nil
}

func CertConfig(certificate, certificateKey string) (*watch.Watch, error) {
	cmWatcher, err := watch.New(logger.Log)
	if err != nil {
		return nil, err
	}
	cm, err := cert.New(certificate, certificateKey)
	if err == watch.ErrFileNotFound {
		return nil, err
	} else {
		if err != nil {
			return nil, err
		} else {
			cmWatcher.Watching = cm
			if err := cmWatcher.Watch(); err != nil {
				return nil, err
			}
		}
	}
	return cmWatcher, nil
}
