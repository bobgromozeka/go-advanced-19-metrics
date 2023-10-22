package agent

// StartupConfig agent startup configuration
type StartupConfig struct {
	ServerAddr     string
	ServerScheme   string
	HashKey        string
	PollInterval   int
	ReportInterval int
	PublicKeyPath  string
}
