package config

import (
	"errors"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/erizzardi/vault-secrets-operator/pkg/utils"
)

type Config struct {
	// Vault config
	VaultUrl   string `config:"vault-url"`
	VaultToken string `config:"vault-token"`

	// Logging configuration
	LogLevel log.Level
}

func (c Config) ValidateConfiguration() error {
	if c.VaultToken == defaultVaultToken {
		return errors.New("Vault token not set.")
	}
	return nil
}

func GetConfigOrDie() (Config, error) {

	vaultUrl := getStringOption(vaultUrlConfig, defaultVaultUrl, fmt.Sprintf("Hashicorp Vault URL - defaust %s", defaultVaultUrl))
	vaultToken := getStringOption(tokenConfig, defaultVaultToken, "Hashicorp Vault token")
	logLevel := getStringOption(logLevelConfig, defaultLogLevel, fmt.Sprintf("Log level - default %s", defaultLogLevel))
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return Config{}, err
	}

	return Config{
		VaultUrl:   *vaultUrl,
		VaultToken: *vaultToken,
		LogLevel:   level}, nil
}

// getStringOption looks through flags and env variables to find the "name" option, and returns a string
func getStringOption(name string, defaultValue string, usage string) *string {
	return flag.String(name, utils.GetEnvOrFallback(utils.FlagToEnv(name), defaultValue), usage)
}
