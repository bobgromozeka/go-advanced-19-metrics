package agent

type StartupConfig struct {
	ServerAddr     string
	ServerScheme   string
	PollInterval   int
	ReportInterval int
	HashKey        string
}
