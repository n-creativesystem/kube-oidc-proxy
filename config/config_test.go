package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/n-creativesystem/oidc-proxy/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	example := config.Config{
		Logging: config.Logging{
			Level: "debug or info or warn or warning(warn) or error or err(error) or critical or dev(debug) or prod(info)",
		},
		Port:              8080,
		SslCertificate:    "ssl/sever.crt",
		SslCertificateKey: "ssl/sever.key",
		Servers: []*config.Servers{
			{
				Login:      "/oauth2/login",
				Callback:   "/oauth2/callback",
				Logout:     "/oauth2/logout",
				ServerName: "virtual sever name",
				Logging: config.Logging{
					Level: "debug or info or warn or warning(warn) or error or err(error) or critical or dev(debug) or prod(info)",
				},
				Oidc: config.Oidc{
					Scopes:       []string{"email", "openid", "offline_access", "profile"},
					Provider:     "https://keycloak/",
					ClientId:     "xxx",
					ClientSecret: "xxx",
					Logout:       "https://keycloak/logout?returnTo=http://localhost:8080/oauth2/login",
					RedirectUrl:  "http://localhost:8080/oauth2/callback",
				},
				Locations: []config.Locations{
					{
						ProxyPass: "http://localhost",
						Urls: []config.Urls{
							{
								Path:  "/",
								Token: "id_token",
							},
						},
					},
				},
				CacheConfig: config.Cache{
					Name:      "memory or etcd",
					Codecs:    []string{},
					Endpoints: []string{},
					CacheTime: 30,
				},
			},
		},
	}
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{}
	exts := []string{".yaml", ".yml", ".json", ".toml"}
	for _, ext := range exts {
		filename := "test" + ext
		os.Remove(filename)
		tests = append(tests, struct {
			name string
			fn   func(t *testing.T)
		}{
			name: fmt.Sprintf("write config to %s", filename),
			fn: func(t *testing.T) {
				err := example.Output(filename)
				assert.NoError(t, err)
				isExists := fileIsExists(filename)
				assert.Equal(t, true, isExists)
			},
		})
		tests = append(tests, struct {
			name string
			fn   func(t *testing.T)
		}{
			name: fmt.Sprintf("read config of %s", filename),
			fn: func(t *testing.T) {
				filename := "test.yaml"
				conf, err := config.ReadConfig(filename)
				if assert.NoError(t, err) {
					is := reflect.DeepEqual(conf, example)
					assert.Equal(t, true, is)
				}
			},
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
	// for _, ext := range exts {
	// 	filename := "test" + ext
	// 	os.Remove(filename)
	// }
}

func fileIsExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
