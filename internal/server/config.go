package server

type StartupConfig struct {
	ServerAddr      string
	FileStoragePath string
	DatabaseDsn     string
	HashKey         string
	StoreInterval   uint
	Restore         bool
	PrivateKeyPath  string
}
