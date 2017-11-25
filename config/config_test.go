package config_test

import (
	"testing"

	"github.com/ifreddyrondon/gobastion/config"
)

func TestNewConfig(t *testing.T) {
	cfg := config.NewConfig()
	if cfg.API.BasePath != "/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
}

func TestLoadFromJSONFile(t *testing.T) {
	cfg := config.NewConfig()
	cfg.FromFile("./testdata/config_test.json")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	cfg := config.NewConfig()
	cfg.FromFile("./testdata/config_test.yaml")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
}
