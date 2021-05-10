#!/bin/env bash

# openssl genrsa 2048 > ssl/server.key

# openssl req -new -key ssl/server.key -subj "/C=JP/ST=Nara/L=Kashiba/O=n-creativesystem/CN=proxy.n-creativesystem.com" > ssl/server.csr

# openssl x509 -days 3650 -req -sha256 -signkey ssl/server.key < ssl/server.csr > ssl/server.crt

openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -subj "/C=JP/ST=Nara/L=Kashiba/O=n-creativesystem/CN=proxy.n-creativesystem.com"
