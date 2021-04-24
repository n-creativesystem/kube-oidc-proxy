# kube-oidc-proxy

kubernetes dashboard token

# Overview

# Environment variables

| Key                | Overview                                              |
| :----------------- | :---------------------------------------------------- |
| OIDC_PROVIDER      | Issur                                                 |
| OIDC_CLIENT_ID     | idp client_id                                         |
| OIDC_CLIENT_SECRET | idp client_secret                                     |
| OIDC_REDIRECT_URL  | redirect callback url                                 |
| OIDC_SCOPE         | scope [email, offline_access, profile, openid etc...] |
| OIDC_LOGOUT        | idp logout url                                        |
| K8S_DASHBOARD      | kubernetes dashbord url                               |
| PROXY_PORT         | listen port number                                    |
| SSL_CERT_FILE      | cert file                                             |
| SSL_KEY_FILE       | key file                                              |
| LOG_LEVEL          | log output level                                      |
