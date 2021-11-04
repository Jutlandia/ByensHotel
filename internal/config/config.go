package config

import (
	"fmt"
	"os"
)

// Config contains the configurations of ByensHotel.
type Config struct {
	Server Server
}

// Server contains the configuration of the server.
type Server struct {
	Port    string
	CsrfKey string
}

// New initializes a new Config.
func New() (*Config, error) {
	csrfKey := os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		return nil, fmt.Errorf("missing csrf key")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}
	cfg := &Config{
		Server: Server{
			Port:    port,
			CsrfKey: csrfKey,
		},
	}
	return cfg, nil
}
