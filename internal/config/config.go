package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config contains the configurations of ByensHotel.
type Config struct {
	SessionKey string
	Server     Server
	LDAP       LDAP
}

// Server contains the configuration of the server.
type Server struct {
	Port    int
	CsrfKey string
}

type LDAP struct {
	URL          string
	BindUsername string
	BindPassword string
	BaseDN       string
}

// New initializes a new Config.
func New() (*Config, error) {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		return nil, fmt.Errorf("mission session key")
	}
	csrfKey := os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		return nil, fmt.Errorf("missing csrf key")
	}
	port, err := getIntEnv("PORT", 1337, "invalid port")
	if err != nil {
		return nil, err
	}
	ldapURL := os.Getenv("LDAP_URL")
	if ldapURL == "" {
		return nil, fmt.Errorf("missing ldap url")
	}
	bindUsername := os.Getenv("LDAP_BIND_USERNAME")
	if bindUsername == "" {
		return nil, fmt.Errorf("missing ldap bind username")
	}
	bindPassword := os.Getenv("LDAP_BIND_PASSWORD")
	if bindPassword == "" {
		return nil, fmt.Errorf("missing ldap bind password")
	}
	baseDN := os.Getenv("LDAP_BASE_DN")
	if baseDN == "" {
		return nil, fmt.Errorf("missing ldap base DN")
	}
	cfg := &Config{
		SessionKey: sessionKey,
		Server: Server{
			Port:    port,
			CsrfKey: csrfKey,
		},
		LDAP: LDAP{
			URL:          ldapURL,
			BindUsername: bindUsername,
			BindPassword: bindPassword,
			BaseDN:       baseDN,
		},
	}
	return cfg, nil
}

func getIntEnv(key string, d int, errMsg string) (int, error) {
	k := os.Getenv(key)
	if k == "" {
		return d, nil
	}
	intEnv, err := strconv.Atoi(k)
	if err != nil {
		return -1, fmt.Errorf("%s: %s", errMsg, k)
	}
	return intEnv, nil
}
