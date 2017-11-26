package config_test

import (
	"testing"

	"github.com/ifreddyrondon/gobastion/config"
)

func TestNewConfig(t *testing.T) {
	cfg := config.New()
	if cfg.API.BasePath != "/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":8080" {
		t.Errorf("Expected BasePath to be ':8080'. Got '%v'", cfg.Server.Addr)
	}
}

func TestLoadFromJSONFile(t *testing.T) {
	cfg := config.New()
	cfg.FromFile("./testdata/config_test.json")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":3000" {
		t.Errorf("Expected BasePath to be ':3000'. Got '%v'", cfg.Server.Addr)
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	cfg := config.New()
	cfg.FromFile("./testdata/config_test.yaml")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":3000" {
		t.Errorf("Expected BasePath to be ':3000'. Got '%v'", cfg.Server.Addr)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg := config.New()
	err := cfg.FromFile("a.yaml")
	if err.Error() != config.ErrorMissingConfigFile.Error() {
		t.Fatalf("Expected Error to be '%v'. Got '%v'", config.ErrorMissingConfigFile.Error(), err)
	}
}

func TestUnmarshalFile(t *testing.T) {
	cfg := config.New()
	err := cfg.FromFile("./testdata/bad_json_test.json")
	if err.Error() != config.ErrorUnmarshalConfig.Error() {
		t.Fatalf("Expected Error to be '%v'. Got '%v'", config.ErrorUnmarshalConfig.Error(), err)
	}
}
