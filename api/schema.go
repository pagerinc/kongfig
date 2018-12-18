package api

type Config struct {
	Host     string    `json:"host"`
	HTTPS    bool      `json:"https"`
	version  string    `json:"version"`
	Services []Service `json:"services,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	CustomID string `json:"custom_id"`
}

type Service struct {
	ConnectTimeout int     `json:"connect_timeout,omitempty"`
	Name           string  `json:"name,omitempty"`
	Path           string  `json:"path,omitempty"`
	ReadTimeout    int     `json:"read_timeout,omitempty"`
	Retries        int     `json:"retries,omitempty"`
	URL            string  `json:"url,omitempty"`
	WriteTimeout   int     `json:"write_timeout,omitempty"`
	Routes         []Route `json:"routes,omitempty"`
	Plugin         Plugin  `json:"plugin,omitempty"`
}

type Routes struct {
	Next string  `json:"next"`
	Data []Route `json:"data"`
}

type Route struct {
	ID           string   `json:"id,omitempty"`
	Protocols    []string `json:"protocols,omitempty"`
	Methods      []string `json:"methods,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	Paths        []string `json:"paths,omitempty"`
	StripPath    bool     `json:"strip_path,omitempty"`
	PreserveHost bool     `json:"preserve_host,omitempty"`
	Service      string   `json:"-"`
}

type Plugin struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	PluginConfig `json:"config,omitempty"`
}

type PluginConfig struct {
	ClaimsToVerify     string `json:"claims_to_verify,omitempty"`
	URIParamNames      string `json:"uri_param_names,omitempty"`
	Credentials        bool   `json:"credentials,omitempty"`
	Origins            string `json:"origins,omitempty"`
	PreflightContinue  bool   `json:"preflight_continue,omitempty"`
	ExposedHeaders     string `json:"exposed_headers,omitempty"`
	Headers            string `json:"headers,omitempty"`
	EchoDownstream     bool   `json:"echo_downstream,omitempty"`
	HeaderName         string `json:"header_name,omitempty"`
	Generator          string `json:"generator,omitempty"`
	Policy             string `json:"policy,omitempty"`
	Hour               int    `json:"hour,omitempty"`
	Second             int    `json:"second,omitempty"`
	AllowedPayloadSize int    `json:"allowed_payload_size,omitempty"`
}
