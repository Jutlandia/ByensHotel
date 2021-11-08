package config_test

import (
	"log"
	"os"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/config"
)

func TestNew(t *testing.T) {
	configs := map[string]string{
		"SESSION_KEY":        "some-32-byte-key",
		"PORT":               "9999",
		"CSRF_KEY":           "some-32-byte-key",
		"LDAP_URL":           "ldap://localhost:10389",
		"LDAP_BIND_USERNAME": "cn=admin,dc=planetexpress,dc=com",
		"LDAP_BIND_PASSWORD": "GoodNewsEveryone",
		"LDAP_BASE_DN":       "dc=planetexpress,dc=com",
	}
	for k, v := range configs {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	expected := struct {
		SessionKey       string
		Port             int
		CsrfKey          string
		URL              string
		LDAPBindUsername string
		LDAPBindPassword string
		LDAPBaseDN       string
	}{
		SessionKey:       "some-32-byte-key",
		Port:             9999,
		CsrfKey:          "some-32-byte-key",
		URL:              "ldap://localhost:10389",
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
	if cfg.LDAP.URL != expected.URL {
		t.Errorf("Expected LDAP url: %s\nGot: %s\n",
			expected.URL, cfg.LDAP.URL)
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
