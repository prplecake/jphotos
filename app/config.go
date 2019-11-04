package app

// A Configuration is the main config object.
type Configuration struct {
	App       Config
	DB        DatabaseConfig `yaml:"database"`
	Templates TemplateConfig
	Uploads   UploadConfig
}

// A Config holds app-specific configuration.
type Config struct {
	Port string
}

// A DatabaseConfig holds database-specific configuration.
type DatabaseConfig struct {
	Username, Password, Hostname, Name string
	Port                               int
}

// A TemplateConfig holds template-specific configuration.
type TemplateConfig struct {
	Path string
}

// An UploadConfig holds upload-specific configuration.
type UploadConfig struct {
	Path           string
	ThumbnailsPath string
}
