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
		Address:       DefaultServerAddress,
		Port:          DefaultServerPort,
		DownloadPath:  DefaultDownloadPath,
		ServeFrontend: true,
	}

	serverCorsConfig := ServerCorsConfig{
		AllowedOrigins: []string{"*"},
	}

	return Config{
		Logger:        loggerConfig,
		Server:        serverConfig,
		Cors:          serverCorsConfig,
		ObjectStorage: ObjectStorageConfig{},
	}
}
