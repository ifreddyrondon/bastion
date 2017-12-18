package config_test

import (
	"testing"

	"github.com/ifreddyrondon/gobastion/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.API.BasePath != "/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":8080" {
		t.Errorf("Expected Addr to be ':8080'. Got '%v'", cfg.Server.Addr)
	}
	if cfg.Debug {
		t.Errorf("Expected Debug to be 'true'. Got '%v'", cfg.Debug)
	}
}

func TestLoadFromJSONFile(t *testing.T) {
	cfg, _ := config.FromFile("./testdata/config_test.json")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":3000" {
		t.Errorf("Expected Addr to be ':3000'. Got '%v'", cfg.Server.Addr)
	}
	if !cfg.Debug {
		t.Errorf("Expected Debug to be 'true'. Got '%v'", cfg.Debug)
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	cfg, _ := config.FromFile("./testdata/config_test.yaml")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":3000" {
		t.Errorf("Expected Addr to be ':3000'. Got '%v'", cfg.Server.Addr)
	}
	if !cfg.Debug {
		t.Errorf("Expected Debug to be 'true'. Got '%v'", cfg.Debug)
	}
}

func TestLoadFromPartialYAMLFile(t *testing.T) {
	cfg, _ := config.FromFile("./testdata/partial_config_test.yaml")
	if cfg.API.BasePath != "/api/" {
		t.Errorf("Expected BasePath to be '/api/'. Got '%v'", cfg.API.BasePath)
	}
	if cfg.Server.Addr != ":8080" {
		t.Errorf("Expected Addr to be ':8080'. Got '%v'", cfg.Server.Addr)
	}
	if cfg.Debug {
		t.Errorf("Expected Debug to be 'false'. Got '%v'", cfg.Debug)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.FromFile("a.yaml")
	if err.Error() != config.ErrorMissingConfigFile.Error() {
		t.Fatalf("Expected Error to be '%v'. Got '%v'", config.ErrorMissingConfigFile.Error(), err)
	}
}

func TestUnmarshalFile(t *testing.T) {
	_, err := config.FromFile("./testdata/bad_json_test.json")
	if err.Error() != config.ErrorUnmarshalConfig.Error() {
		t.Fatalf("Expected Error to be '%v'. Got '%v'", config.ErrorUnmarshalConfig.Error(), err)
	}
}
