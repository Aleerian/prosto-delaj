package models

type AppState struct {
	Env           *Environment
	ConfigService *ConfigService
	ConfigVault   *ConfigVault
}

type Environment struct {
	VaultAddr    string `env:"VAULT_ADDR"`
	VaultRoleID  string `env:"VAULT_ROLE_ID"`
	IsReadConfig bool   `env:"IS_READ_CONFIG_FILE"`
}

type ConfigService struct {
	Server     *ServerConfig     `json:"server" binding:"required"`
	BusinessDB *BusinessDBConfig `json:"business-database" binding:"required"`
}

type ServerConfig struct {
	Port       string `json:"server_port" binding:"required"`
	ServerMode string `json:"server_mode" binding:"required"`
	Domain     string `json:"server_domain" binding:"required"`
}

type BusinessDBConfig struct {
	Password string `json:"db_password" binding:"required"`
	Host     string `json:"db_host" binding:"required"`
	Port     string `json:"db_port" binding:"required"`
	Username string `json:"db_username" binding:"required"`
	DBName   string `json:"db_name" binding:"required"`
	SSLMode  string `json:"db_ssl_mode" binding:"required"`
}

type VaultPathConfig struct {
	Name   string `yaml:"name"`
	Engine string `yaml:"engine"`
	Path   string `yaml:"path"`
	Field  string `yaml:"field"`
}

type ConfigVault struct {
	WrappedSecretID    string `json:"auth_token" binding:"required"`
	WrappedSecretPaths string `json:"wrapping_token" binding:"required"`
}

type VaultSecretsConfig struct {
	Secrets VaultPathConfig `yaml:"secret"`
}
