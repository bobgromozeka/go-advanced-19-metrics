package server

type StartupConfig struct {
	HTTPAddr           string `json:"http_address"`
	GRPCAddr           string `json:"grpc_address"`
	FileStoragePath    string `json:"store_file"`
	DatabaseDsn        string `json:"database_dsn"`
	HashKey            string `json:"hash_key"`
	StoreInterval      uint   `json:"store_interval"`
	Restore            bool   `json:"restore"`
	PrivateKeyPath     string `json:"crypto_key"`
	TrustedSubnet      string `json:"trusted_subnet"`
	GRPCPrivateKeyPath string `json:"grpc_private_key_path"`
	GRPCCertPath       string `json:"grpc_cert_path"`
}
