package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/erizzardi/vault-secrets-operator/pkg/utils"
)

type Config struct {
	// Vault config
	VaultUrl   string `yaml:"vault-url"`
	VaultToken string `yaml:"vault-token"`

	// Controller config
	ResyncPeriod int `yaml:"resyncPeriod"`
	LoopPeriod   int `yaml:"loopPeriod"`
	LogLevel     log.Level

	// Kubernetes config
	LocalTesting bool   `yaml:"kubeconfig"`
	Kubeconfig   string `yaml:"kubeconfig"`
}

func (c Config) ValidateConfiguration() error {
	if c.VaultToken == noDefaultString {
		return errors.New("Vault token not set")
	}
	return nil
}

func GetConfigOrDie() (Config, error) {

	vaultUrl := getStringOption(vaultUrlConfig, defaultVaultUrl, "Hashicorp Vault URL")
	vaultToken := getStringOption(tokenConfig, noDefaultString, "Hashicorp Vault token")
	resyncPeriod, err := getIntOption(resyncPeriodConfig, defaultResyncPeriod, "Controller Informer resync period, in seconds")
	if err != nil {
		return Config{}, errors.New(err.Error() + ": resync-period")
	}
	loopPeriod, err := getIntOption(loopPeriodConfig, defaultLoopPeriod, "Controller main loop period, in seconds")
	if err != nil {
		return Config{}, errors.New(err.Error() + ": loop-period")
	}
	logLevel := getStringOption(logLevelConfig, defaultLogLevel, "Select the logging level")
	localTesting := getBoolOption(localTestingConfig, defaultLocalTesting, "Toggle to enable local testing mode. Specify the kubeconfig path with --kubeconfig")
	kubeconfig := getStringOption(kubeconfigConfig, defaultKubeconfig, "Location of the kubeconfig file, ignored if --local-testing is not set")
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return Config{}, err
	}

	return Config{
		VaultUrl:     *vaultUrl,
		VaultToken:   *vaultToken,
		ResyncPeriod: *resyncPeriod,
		LoopPeriod:   *loopPeriod,
		LogLevel:     level,
		LocalTesting: *localTesting,
		Kubeconfig:   *kubeconfig}, nil
}

// getStringOption looks through flags and env variables to find the "name" option, and returns a string
func getStringOption(name string, defaultValue string, usage string) *string {
	return flag.String(name, utils.GetEnvOrFallback(utils.FlagToEnv(name), defaultValue), usage)
}

// getIntOption looks through flags and env variables to find the "name" option, and returns an int
func getIntOption(name string, defaultValue int, usage string) (*int, error) {
	def, err := strconv.Atoi(utils.GetEnvOrFallback(utils.FlagToEnv(name), fmt.Sprintf("%d", defaultValue)))
	if err != nil {
		return nil, err
	}
	return flag.Int(name, def, usage), nil
}

// getBoolOption looks through flags and env variables to find the "name" option, and returns a bool
func getBoolOption(name string, defaultValue bool, usage string) *bool {
	return flag.Bool(name, defaultValue, usage)
}
