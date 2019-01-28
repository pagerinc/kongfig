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

// Route represents a route for a microservice
type Route struct {
	Name          string   `yaml:"name,omitempty" json:"-"`
	ID            string   `yaml:"id,omitempty" json:"id,omitempty"`
	Service       string   `yaml:"apply_to,omitempty" json:"service,omitempty"`
	Hosts         []string `yaml:"hosts,omitempty" json:"hosts,omitempty"`
	Paths         []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	Methods       []string `yaml:"methods,omitempty" json:"methods,omitempty"`
	StripPath     bool     `yaml:"strip_path,omitempty" json:"strip_path,omitempty"`
	Protocols     []string `yaml:"protocols,omitempty" json:"protocols,omitempty"`
	RegexPriority int      `yaml:"regex_priority,omitempty" json:"regex_priority,omitempty"`
	PreserveHost  bool     `yaml:"preserve_host,omitempty" json:"preserve_host,omitempty"`
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
	Protocol       string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

// Services represents the response body returned from GET /services, a Kong API endpoint
// Contains service data for all services
type Services struct {
	Next string    `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Service `yaml:"data,omitempty" json:"data,omitempty"`
}

// Routes represents the response body returned from GET /routes, a Kong API endpoint
// Contains route data for all routes
type Routes struct {
	Next string  `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Route `yaml:"data,omitempty" json:"data,omitempty"`
}

// Consumer represents the user credential for authentication to Kong
type Consumer struct {
	Username string `json:"username" yaml:"username"`
	CustomID string `json:"custom_id,omitempty" yaml:"custom_id"`
}

// Consumers represents the response body returned from GET /consumers, a Kong API endpoint
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

// Plugins represents the response body of GET /plugins endpoint of Kong Admin API
type Plugins struct {
	Next string   `yaml:"next,omitempty" json:"next,omitempty"`
	Data []Plugin `yaml:"data,omitempty" json:"data,omitempty"`
}

// Plugin represents a feature or middleware in Kong
type Plugin struct {
	ID       string                 `yaml:"id,omitempty" json:"id,omitempty"`
	Name     string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Enabled  bool                   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Services []string               `yaml:"services,omitempty" json:"-"`
	Routes   []string               `yaml:"routes,omitempty" json:"-"`
	Target   string                 `yaml:"target,omitempty" json:"-"`
	Config   map[string]interface{} `yaml:"config,omitempty" json:"config,omitempty"`
}

type HeaderList struct {
	Headers []string `yaml:"headers,omitempty" json:"headers,omitempty"`
}
