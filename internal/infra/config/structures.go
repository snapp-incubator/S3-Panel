package config

type LoggerConfig struct {
	Level string `json:"level" koanf:"level"`
}

type ServerConfig struct {
	Address string `json:"address" koanf:"address"`
	Port    string `json:"port"    koanf:"port"`
}

type ServerCorsConfig struct {
	AllowedOrigins []string `json:"allowed_origins" koanf:"allowed_origins"`
}

type ObjectStorageConfig struct {
	URL            string `json:"url"         koanf:"url"`
	AccessKeyAdmin string `json:"access_key"  koanf:"access_key"`
	SecretKeyAdmin string `json:"secret_key"  koanf:"secret_key"`
}
