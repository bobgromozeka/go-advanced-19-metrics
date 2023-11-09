package agent

// StartupConfig agent startup configuration
type StartupConfig struct {
	ServerAddr     string `json:"address"`
	ServerScheme   string `json:"scheme"`
	HashKey        string `json:"hash_key"`
	PollInterval   int    `json:"poll_interval"`
	ReportInterval int    `json:"report_interval"`
	PublicKeyPath  string `json:"crypto_key"`
	ReportGRPC     bool   `json:"report_grpc"`
}
