package config

import (
	"net/http"

	"github.com/namsral/flag"
	"gitlab.com/metronero/backend/pkg/api"
)

var (
	Hostname  string
	Bind      string
	Debug     bool
	Backend   string
	Uname     string
	JwtSecret string
	Api       *api.ApiClient
)

func Load() {
	flag.StringVar(&Hostname, "hostname", "http://localhost:5002", "User facing hostname or domain name with protocol scheme")
	flag.StringVar(&Bind, "bind", ":5002", "Bind address")
	flag.BoolVar(&Debug, "debug", false, "Debug mode")
	flag.StringVar(&Backend, "backend", "http://localhost:5001", "Metronero backend host:port")
	flag.StringVar(&JwtSecret, "token-secret", "secret", "Secret for authentication tokens, same as backend")
	Api = &api.ApiClient{
		Client:  http.DefaultClient,
		BaseUrl: Backend,
	}
	flag.Parse()
}
