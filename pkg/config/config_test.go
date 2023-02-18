package config

import (
	"testing"

	"github.com/sirupsen/logrus"
)

var validConfig = Config{
	VaultUrl:     defaultVaultUrl,
	VaultToken:   "my-token",
	ResyncPeriod: defaultResyncPeriod,
	LoopPeriod:   defaultLoopPeriod,
	LogLevel:     logrus.InfoLevel,
	LocalTesting: defaultLocalTesting,
	Kubeconfig:   defaultKubeconfig,
}
var invalidConfig1 = Config{
	VaultToken: "",
}
var invalidConfig2 = Config{
	LoopPeriod: -1,
}
var invalidConfig3 = Config{
	ResyncPeriod: -1,
}

func TestValidConfig(t *testing.T) {
	if err := validConfig.ValidateConfiguration(); err != nil {
		t.Error("Valid configuration recognized as invalid")
	}
}

func TestInvalidConfig(t *testing.T) {
	if err := invalidConfig1.ValidateConfiguration(); err == nil {
		t.Error("Invalid configuration recognized as valid: unset vault token")
	}
	if err := invalidConfig2.ValidateConfiguration(); err == nil {
		t.Error("Invalid configuration recognized as valid: LoopPeriod negative")
	}
	if err := invalidConfig3.ValidateConfiguration(); err == nil {
		t.Error("Invalid configuration recognized as valid: ResyncPeriod negative")
	}
}
