package config

// Default values
const (
	// Vault configuration
	defaultVaultUrl   = "http://localhost:8200"
	defaultVaultToken = ""
	// Logging configuration
	defaultLogLevel = "INFO"
)

// Miscellanea
const (
	GroupName               = "erizzardi.mine.io"
	SecretVersionAnnotation = GroupName + "/secret-version"
	ManagedAnnotation       = GroupName + "/managed"
	LACAnnotation           = GroupName + "/last-applied-configuration"
)
