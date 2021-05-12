package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/n-creativesystem/oidc-proxy/app"
	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/plugin"
	"github.com/n-creativesystem/oidc-proxy/session"
	"github.com/n-creativesystem/oidc-proxy/store"
	"github.com/prometheus/common/log"
	"golang.org/x/oauth2"

	toml "github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	Logging           Logging             `yaml:"logging" toml:"logging" json:"logging"`
	Servers           []*Servers          `yaml:"servers" toml:"servers" json:"servers"`
	Port              int                 `yaml:"port" toml:"port" json:"port"`
	SslCertificate    string              `yaml:"ssl_certificate" toml:"ssl_certificate" json:"ssl_certificate"`
	SslCertificateKey string              `yaml:"ssl_certificate_key" toml:"ssl_certificate_key" json:"ssl_certificate_key"`
	port              string              `yaml:"-" toml:"-" json:"-"`
	mapSrvConfig      map[string]*Servers `yaml:"-" toml:"-" json:"-"`
}

func (c *Config) GetPort() string {
	if c.port == "" {
		port := strconv.Itoa(c.Port)
		if port != "" && port[:1] != ":" {
			port = ":" + port
		}
		c.port = port
	}
	return c.port
}

func (c *Config) GetServerConfig(serverName string) *Servers {
	return c.mapSrvConfig[serverName]
}

func (c *Config) Output(filename string) error {
	typ := GetExtension(filename)
	var marshal func(interface{}) ([]byte, error)
	switch typ {
	case Yaml:
		marshal = yaml.Marshal
	case Toml:
		marshal = toml.Marshal
	case Json:
		marshal = json.Marshal
	}
	buf, err := marshal(c)
	if err != nil {
		return err
	}
	if typ == Json {
		var b bytes.Buffer
		if err := json.Indent(&b, buf, "", " "); err != nil {
			return err
		}
		buf = b.Bytes()
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(buf)
	return err
}

// Servers
type Servers struct {
	Oidc       Oidc        `yaml:"oidc" toml:"oidc" json:"oidc"`
	Locations  []Locations `yaml:"locations" toml:"locations" json:"locations"`
	Logging    Logging     `yaml:"logging" toml:"logging" json:"logging"`
	CookieName string      `yaml:"cookie_name" toml:"cookie_name" json:"cookie_name"`
	ServerName string      `yaml:"server_name" toml:"server_name" json:"server_name"`
	Session    Session     `yaml:"session" toml:"session" json:"session"`
	Login      string      `yaml:"login" toml:"login" json:"login"`
	Callback   string      `yaml:"callback" toml:"callback" json:"callback"`
	Logout     string      `yaml:"logout" toml:"logout" json:"logout"`
	Redirect   bool        `yaml:"redirect" toml:"redirect" json:"redirect"`
	port       string      `yaml:"-" toml:"-" json:"-"`
}

func (s *Servers) newSessionClient(client *hplugin.Client) session.Session {
	var storage session.Session
	rpcClient, err := client.Client()
	if err != nil {
		log.Error(err)
	}
	if rpcClient != nil {
		raw, err := rpcClient.Dispense("session")
		if err != nil {
			log.Error(err)
		}
		if raw != nil {
			if val, ok := raw.(*plugin.GRPCSessionClient); ok {
				val.PluginClient = client
				storage = val
			}
		}
	}
	// if storage == nil {
	// 	storage = session.New()
	// }
	storage.Init(context.TODO(), s.Session.Args)
	return storage
}

func (s *Servers) init(port string) {
	cmd := exec.Command(s.Session.GetSessionPlugin())
	var sessionPlugin *hplugin.Client
	var dispose *app.Dispose
	dispose = app.Store.Dispose(s.ServerName)
	if dispose != nil && dispose.Plugin != nil {
		dispose.Close()
	}
	sessionPlugin = hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig:  plugin.Handshake,
		VersionedPlugins: plugin.VersionedPlugins,
		Cmd:              cmd,
		AllowedProtocols: []hplugin.Protocol{hplugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.LevelFromString(s.Logging.Level),
			Output:     s.Logging.writer,
			TimeFormat: "2006/01/02 - 15:04:05",
			Name:       "session",
		}),
		Managed: true,
	})
	storage := s.newSessionClient(sessionPlugin)
	codecs := s.Session.GetCodecs()
	store := store.NewStore(storage, codecs...)
	app.Store.Add(s.ServerName, &app.Dispose{
		Store:  store,
		Plugin: sessionPlugin,
	})
	s.port = port
}

func (s *Servers) Is() error {
	msg := func(message string) string {
		return fmt.Sprintf("%s: %s", s.ServerName, message)
	}
	if !s.Session.IsCodecs() {
		return errors.New(msg("no codecs provided"))
	}

	return nil
}

func (s *Servers) GetHostname() string {
	return s.ServerName + s.port
}

// Oidc
type Oidc struct {
	Scopes       []string `yaml:"scopes" toml:"scopes" json:"scopes"`
	Provider     string   `yaml:"provider" toml:"provider" json:"provider"`
	ClientId     string   `yaml:"client_id" toml:"client_id" json:"client_id"`
	ClientSecret string   `yaml:"client_secret" toml:"client_secret" json:"client_secret"`
	RedirectUrl  string   `yaml:"redirect_url" toml:"redirect_url" json:"redirect_url"`
	Logout       string   `yaml:"logout" toml:"logout" json:"logout"`
	// GrantType    string    `yaml:"grant_type" toml:"grant_type" json:"grant_type"`
	Audiences []string `yaml:"audiences" toml:"audiences" json:"audiences"`
}

func (o *Oidc) SetValues() []oauth2.AuthCodeOption {
	var authCodeOptions []oauth2.AuthCodeOption
	var audiences Audiences
	for _, audience := range o.Audiences {
		audiences = append(audiences, Audience(audience))
	}
	authCodeOptions = append(authCodeOptions, audiences.SetValue()...)
	return authCodeOptions
}

// Audience
type Audience string

func (a Audience) String() string {
	return string(a)
}

type Audiences []Audience

func (a Audiences) SetValue() []oauth2.AuthCodeOption {
	var audiences []oauth2.AuthCodeOption
	for _, audience := range a {
		audiences = append(audiences, oauth2.SetAuthURLParam("audience", audience.String()))
	}
	return audiences
}

// Session
type Session struct {
	Name   string                 `yaml:"name" toml:"name" json:"name"`
	Codecs []string               `yaml:"codecs" toml:"codecs" json:"codecs"`
	Args   map[string]interface{} `yaml:"args" toml:"args" json:"args"`
}

const defaultSessionPlugin = "memory"

func (c *Session) GetSessionPlugin() string {
	const pluginDir = "oidc-plugin"
	if c.Name != "" {
		return filepath.Join("./", pluginDir, c.Name)
	}
	return filepath.Join("./", pluginDir, defaultSessionPlugin)
}

func (c *Session) IsCodecs() bool {
	return len(c.Codecs) > 0
}

func (c *Session) GetCodecs() [][]byte {
	var codecs [][]byte
	for _, codec := range c.Codecs {
		codecs = append(codecs, []byte(codec))
	}
	return codecs
}

// Locations
type Locations struct {
	ProxyPass      string `yaml:"proxy_pass" toml:"proxy_pass" json:"proxy_pass"`
	ProxySSLVerify string `yaml:"proxy_ssl_verify" toml:"proxy_ssl_verify" json:"proxy_ssl_verify"`
	Urls           []Urls `yaml:"urls" toml:"urls" json:"urls"`
}

func (l *Locations) IsProxySSLVerify() bool {
	return l.ProxySSLVerify == "on"
}

// Urls
type Urls struct {
	Path  string `yaml:"path" toml:"path" json:"path"`
	Token string `yaml:"token" toml:"token" json:"token"`
	Type  string `yaml:"type" toml:"type" json:"type"`
}

// Logging
type Logging struct {
	Level      string    `yaml:"level" toml:"level" json:"level"`
	FileName   string    `yaml:"filename" toml:"filename" json:"filename"`
	LogFormat  string    `yaml:"logformat" toml:"logformat" json:"logformat"`
	TimeFormat string    `yaml:"timeformat" toml:"timeformat" json:"timeformat"`
	writer     io.Writer `yaml:"-"`
}

func (l *Logging) GetLogger() logger.ILogger {
	writer := []io.Writer{}
	if l.FileName != "" {
		if logfile, err := os.OpenFile(l.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			writer = append(writer, logfile)
		} else {
			logger.Log.Error(err.Error())
		}
	} else {
		writer = append(writer, os.Stdout)
	}
	l.writer = io.MultiWriter(writer...)
	return logger.New(l.writer, logger.Convert(l.Level), logger.ConvertLogFmt(l.LogFormat), logger.ConvertTimeFmt(l.TimeFormat))
}

type GetConfiguration func() Servers

func ReadConfig(filename string) (Config, error) {
	var conf Config
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf, err
	}
	expand := os.ExpandEnv(string(buf))
	ext := GetExtension(filename)
	var unmarshal func([]byte, interface{}) error
	switch ext {
	case Yaml:
		unmarshal = yaml.Unmarshal
	case Toml:
		unmarshal = toml.Unmarshal
	case Json:
		unmarshal = json.Unmarshal
	default:
		unmarshal = yaml.Unmarshal
	}
	err = unmarshal([]byte(expand), &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

type Extension int

const (
	Yaml Extension = iota
	Toml
	Json
)

func GetExtension(filename string) Extension {
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml":
		return Yaml
	case ".toml":
		return Toml
	case ".json":
		return Json
	default:
		return Yaml
	}
}

func New(filename string) (Config, error) {
	conf, err := ReadConfig(filename)
	if err != nil {
		return conf, err
	}
	conf.mapSrvConfig = map[string]*Servers{}
	for _, s := range conf.Servers {
		server := s
		if err := server.Is(); err != nil {
			return conf, err
		}
		server.init(conf.GetPort())
		conf.mapSrvConfig[s.ServerName] = server
	}
	return conf, nil
}
