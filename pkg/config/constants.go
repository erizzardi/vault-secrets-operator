package config

var noDefaultString = ""

// var noDefaultInt = 0

// Default values
const (
	// Vault configuration
	defaultVaultUrl = "http://localhost:8200"

	// Controller configuration
	defaultResyncPeriod = 60
	defaultLoopPeriod   = 1
	defaultLogLevel     = "INFO"

	// Kubernetes configuration
	defaultLocalTesting = false
	defaultKubeconfig   = ".kube/config"
)

// Miscellanea
const (
	GroupName               = "erizzardi.mine.io"
	SecretVersionAnnotation = GroupName + "/secret-version"
	ManagedAnnotation       = GroupName + "/managed"
	LACAnnotation           = GroupName + "/last-applied-configuration"
)
