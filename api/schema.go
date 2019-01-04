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
	Name           string `yaml:"name,omitempty" json:"name,omitempty"`
	URL            string `yaml:"url,omitempty" json:"url,omitempty"`
	Host           string `yaml:"host,omitempty" json:"host,omitempty"`
	Path           string `yaml:"path,omitempty" json:"path,omitempty"`
	Port           int    `yaml:"port,omitempty" json:"port,omitempty"`
	ConnectTimeout int    `yaml:"connect_timeout,omitempty" json:"connect_timeout,omitempty"`
	WriteTimeout   int    `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`
	ReadTimeout    int    `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`
	Retries        int    `yaml:"retries,omitempty" json:"retries,omitempty"`
}

// Route represents routes for each microservice
type Route struct {
	ID            string   `yaml:"id,omitempty" json:"id,omitempty"`
	Service       string   `yaml:"service,omitempty" json:"service,omitempty"`
	Hosts         []string `yaml:"hosts,omitempty" json:"hosts,omitempty"`
	Paths         []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	Methods       []string `yaml:"methods,omitempty" json:"methods,omitempty"`
	StripPath     bool     `yaml:"strip_path,omitempty" json:"strip_path,omitempty"`
	Protocols     []string `yaml:"protocols,omitempty" json:"protocols,omitempty"`
	RegexPriority int      `yaml:"regex_priority,omitempty" json:"regex_priority,omitempty"`
	PreserveHost  bool     `yaml:"preserve_host,omitempty" json:"preserve_host,omitempty"`
}

type Services struct {
	Next string    `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Service `yaml:"data,omitempty" json:"data,omitempty"`
}

type Routes struct {
	Next string  `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Route `yaml:"data,omitempty" json:"data,omitempty"`
}

// Consumer represents the user credential for authentication to Kong
type Consumer struct {
	Username string `yaml:"username"`
	CustomID string `yaml:"custom_id"`
}

type Consumers struct {
	Next string     `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Consumer `yaml:"data,omitempty" json:"data,omitempty"`
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

type Plugins struct {
	Next string   `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Plugin `yaml:"data,omitempty" json:"data,omitempty"`
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
