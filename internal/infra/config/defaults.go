package config

const (
	DefaultLogLevel      = "debug"
	DefaultServerAddress = "127.0.0.1"
	DefaultServerPort    = "8080"
	DefaultDownloadPath  = "/tmp"
)

func DefaultConfig() Config {
	loggerConfig := LoggerConfig{
		Level: DefaultLogLevel,
	}

	serverConfig := ServerConfig{
		Address:      DefaultServerAddress,
		Port:         DefaultServerPort,
		DownloadPath: DefaultDownloadPath,
	}

	serverCorsConfig := ServerCorsConfig{
		AllowedOrigins: []string{"*"},
	}

	return Config{
		LoggerConfigs:        loggerConfig,
		ServerConfigs:        serverConfig,
		ServerCorsConfigs:    serverCorsConfig,
		ObjectStorageConfigs: ObjectStorageConfig{},
	}
}
