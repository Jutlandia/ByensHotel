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
		"PORT":     "9999",
		"CSRF_KEY": "some-32-byte-key",
	}
	loadEnv(configs)
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	expected := struct {
		Port    string
		CsrfKey string
	}{
		Port:    "9999",
		CsrfKey: "some-32-byte-key",
	}
	if cfg.Server.Port != expected.Port {
		t.Errorf("Expected port: %s\nGot: %s\n", expected.Port, cfg.Server.Port)
	}
	if cfg.Server.CsrfKey != expected.CsrfKey {
		t.Errorf("Expected CsrfKey: %s\nGot: %s\n",
			expected.CsrfKey, cfg.Server.CsrfKey)
	}
}
