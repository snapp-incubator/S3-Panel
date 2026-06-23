package config

type LoggerConfig struct {
	Level string `json:"level" koanf:"level"`
}

type ServerConfig struct {
	Address       string `json:"address"         koanf:"address"`
	Port          string `json:"port"            koanf:"port"`
	AuthEnabled   string `json:"auth_enabled"    koanf:"auth_enabled"`
	AuthKeyLookup string `json:"auth_key_lookup" koanf:"auth_key_lookup"`
	AuthToken     string `json:"auth_token"      koanf:"auth_token"`
	DownloadPath  string `json:"download_path"   koanf:"download_path"`
	// ServeFrontend controls whether the embedded SPA is served. Disable it for
	// API-only instances. Defaults to true.
	ServeFrontend bool `json:"serve_frontend"  koanf:"serve_frontend"`
	// Region is the name of the region this instance serves directly (e.g.
	// "teh-1"). Requests whose "region" header matches it (or is empty) are
	// handled locally.
	Region string `json:"region"          koanf:"region"`
	// RegionEndpoints maps other region names to the base URL of their backend.
	// When the frontend is served, a request for one of these regions is
	// reverse-proxied to that backend. Only meaningful on a frontend instance.
	RegionEndpoints map[string]string `json:"region_endpoints" koanf:"region_endpoints"`
}

type ServerCorsConfig struct {
	AllowedOrigins []string `json:"allowed_origins" koanf:"allowed_origins"`
}

type ObjectStorageConfig struct {
	URL            string `json:"url"         koanf:"url"`
	AccessKeyAdmin string `json:"access_key"  koanf:"access_key"`
	SecretKeyAdmin string `json:"secret_key"  koanf:"secret_key"`
}
