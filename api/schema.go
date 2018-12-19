package api

// Config models the top-level structure of the config YAML file
type Config struct {
	Host        string       `yaml:"host"`
	HTTPS       bool         `yaml:"https"`
	Version     string       `yaml:"version"`
	Services    []Service    `yaml:"services"`
	Routes      []Route      `yaml:"routes"`
	Plugins     []Plugin     `yaml:"plugins"`
	Consumers   []Consumer   `yaml:"consumers,omitempty"`
	Credentials []Credential `yaml:"credentials,omitempty"`
}

// Service represents the upstream microservice
type Service struct {
	Name           string `yaml:"name,omitempty"`
	URL            string `yaml:"url,omitempty"`
	Host           string `yaml:"host,omitempty"`
	Path           string `yaml:"path,omitempty"`
	Port           int    `yaml:"port,omitempty"`
	ConnectTimeout int    `yaml:"connect_timeout,omitempty"`
	WriteTimeout   int    `yaml:"write_timeout,omitempty"`
	ReadTimeout    int    `yaml:"read_timeout,omitempty"`
	Retries        int    `yaml:"retries,omitempty"`
}

// Route represents routes for each microservice
type Route struct {
	Name    string      `yaml:"name"`
	ApplyTo string      `yaml:"apply_to,omitempty"`
	Config  RouteConfig `yaml:"config,omitempty"`
}

// RouteConfig represents the config property in Route struct
type RouteConfig struct {
	Hosts     []string `yaml:"hosts"`
	Methods   []string `yaml:"methods,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
	StripPath bool     `yaml:"strip_path,omitempty"`
}

// Consumer represents the user credential for authentication to Kong
type Consumer struct {
	Username string `yaml:"username"`
	CustomID string `yaml:"custom_id"`
}

// Credential represents user
type Credential struct {
	Name   string           `yaml:"name"`
	Target string           `yaml:"target"`
	Config CredentialConfig `yaml:"config"`
}

// CredentialConfig represents the config object inside the Credential struct
type CredentialConfig struct {
	ID     string `yaml:"id"`
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

// Plugin represents a feature or middleware in Kong
type Plugin struct {
	Name    string       `yaml:"name"`
	Enabled bool         `yaml:"enabled,omitempty"`
	Target  []string     `yaml:"target,omitempty"`
	Config  PluginConfig `yaml:"config,omitempty"`
}

// PluginConfig represents the objects in config slive in the Plugin struct
type PluginConfig struct {
	Credentials       bool   `yaml:"credentials,omitempty"`
	Origins           string `yaml:"origins,omitempty"`
	ClaimsToVerify    string `yaml:"claims_to_verify,omitempty"`
	URIParamNames     string `yaml:"uri_param_names,omitempty"`
	PreflightContinue bool   `yaml:"preflight_continue,omitempty"`
	ExposedHeaders    string `yaml:"exposed_headers,omitempty"`
	Headers           string `yaml:"headers,omitempty"`
}
