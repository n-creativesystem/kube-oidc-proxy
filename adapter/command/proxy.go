package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n-creativesystem/oidc-proxy/app"
	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/watch"
	"github.com/n-creativesystem/oidc-proxy/watch/cert"
	watchConfig "github.com/n-creativesystem/oidc-proxy/watch/config"
	"github.com/urfave/cli/v2"
)

const (
	appConf = "config"
	// sessionPlugin = "session-plugin"
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
				return ProxyAction(c.String(appConf))
			},
		},
	},
}

func ProxyAction(configFilename string) error {
	https := true
	appWatcher, err := AppConfig(configFilename)
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
	s := &http.Server{
		Addr:    appConf.GetPort(),
		Handler: multiHost,
	}
	logger.Log.Info("Application Start")
	if https {
		s.TLSConfig = &tls.Config{
			GetCertificate: cmWatcher.Watching.(*cert.Watch).GetCertificate,
		}
		go func() {
			if err := s.ListenAndServeTLS("", ""); err != nil {
				logger.Log.Error(err)
			}
		}()
	} else {
		go func() {
			if err := s.ListenAndServe(); err != nil {
				logger.Log.Error(err)
			}
		}()
	}
	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGSTOP,
	}
	osNotify := make(chan os.Signal, 1)
	signal.Notify(osNotify, signals...)
	sig := <-osNotify
	logger.Log.Info(fmt.Sprintf("signal: %v", sig))
	s.RegisterOnShutdown(func() {
		for _, server := range appConf.Servers {
			if dispose := app.Store.Dispose(server.ServerName); dispose != nil {
				dispose.Close()
			}
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return s.Shutdown(ctx)
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
