package env

type Environments struct {
	OIDC_PROVIDER          string `envconfig:"OIDC_PROVIDER" default:"" required:"true"`
	OIDC_CLIENT_ID         string `envconfig:"OIDC_CLIENT_ID" default:"" required:"true"`
	OIDC_CLIENT_SECRET     string `envconfig:"OIDC_CLIENT_SECRET" default:"" required:"true"`
	OIDC_REDIRECT_URL      string `envconfig:"OIDC_REDIRECT_URL" default:"http://localhost:8080"`
	OIDC_SCOPE             string `envconfig:"OIDC_SCOPE" default:"" required:"true"`
	OIDC_LOGOUT            string `envconfig:"OIDC_LOGOUT" default:"" required:"true"`
	KubernetesDashboardURL string `envconfig:"K8S_DASHBOARD" default:"" required:"true"`
	ProxyPort              string `envconfig:"PROXY_PORT" default:"8080"`
	CertFile               string `envconfig:"SSL_CERT_FILE" default:""`
	KeyFile                string `envconfig:"SSL_KEY_FILE" default:""`
	LogLevel               int    `envconfig:"LOG_LEVEL" default:"3"`
}

var (
	Env Environments
)
