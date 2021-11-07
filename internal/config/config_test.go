package config_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/config"
	"github.com/joho/godotenv"
)

func loadEnv(configs map[string]string) {
	f, err := os.CreateTemp("", ".env.temp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())
	var sb strings.Builder
	for k, v := range configs {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	if _, err := f.Write([]byte(sb.String())); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	godotenv.Load(f.Name())
}

func TestNew(t *testing.T) {
	configs := map[string]string{
		"SESSION_KEY":        "some-32-byte-key",
		"PORT":               "9999",
		"CSRF_KEY":           "some-32-byte-key",
		"LDAP_HOST":          "localhost",
		"LDAP_PORT":          "10389",
		"LDAP_BIND_USERNAME": "cn=admin,dc=planetexpress,dc=com",
		"LDAP_BIND_PASSWORD": "GoodNewsEveryone",
		"LDAP_BASE_DN":       "dc=planetexpress,dc=com",
	}
	loadEnv(configs)
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	expected := struct {
		SessionKey       string
		Port             int
		CsrfKey          string
		LDAPHost         string
		LDAPPort         int
		LDAPBindUsername string
		LDAPBindPassword string
		LDAPBaseDN       string
	}{
		SessionKey:       "some-32-byte-key",
		Port:             9999,
		CsrfKey:          "some-32-byte-key",
		LDAPHost:         "localhost",
		LDAPPort:         10389,
		LDAPBindUsername: "cn=admin,dc=planetexpress,dc=com",
		LDAPBindPassword: "GoodNewsEveryone",
		LDAPBaseDN:       "dc=planetexpress,dc=com",
	}
	if cfg.SessionKey != expected.SessionKey {
		t.Errorf("Expected session key: %s\nGot: %s\n",
			expected.SessionKey, cfg.SessionKey)
	}
	if cfg.Server.Port != expected.Port {
		t.Errorf("Expected port: %d\nGot: %d\n", expected.Port, cfg.Server.Port)
	}
	if cfg.Server.CsrfKey != expected.CsrfKey {
		t.Errorf("Expected CsrfKey: %s\nGot: %s\n",
			expected.CsrfKey, cfg.Server.CsrfKey)
	}
	if cfg.LDAP.Host != expected.LDAPHost {
		t.Errorf("Expected LDAP host: %s\nGot: %s\n",
			expected.LDAPHost, cfg.LDAP.Host)
	}
	if cfg.LDAP.Port != expected.LDAPPort {
		t.Errorf("Expected LDAP port: %d\nGot: %d\n",
			expected.LDAPPort, cfg.LDAP.Port)
	}
	if cfg.LDAP.BindUsername != expected.LDAPBindUsername {
		t.Errorf("Expected LDAP bind username: %s\nGot: %s\n",
			expected.LDAPBindUsername, cfg.LDAP.BindUsername)
	}
	if cfg.LDAP.BindPassword != expected.LDAPBindPassword {
		t.Errorf("Expected LDAP bind password: %s\nGot: %s\n",
			expected.LDAPBindPassword, cfg.LDAP.BindPassword)
	}
	if cfg.LDAP.BaseDN != expected.LDAPBaseDN {
		t.Errorf("Expected LDAP base DN: %s\nGot: %s\n",
			expected.LDAPBaseDN, cfg.LDAP.BaseDN)
	}
}
