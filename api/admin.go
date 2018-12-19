package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	contentType     string = "Content-Type"
	applicationJSON string = "application/json; charset=utf-8"
)

// Client represents the public API
type Client struct {
	config  *Config
	client  *http.Client
	BaseURL string
}

// NewClient returns a Client object with the parsed configuration
func NewClient(filePath string) (*Client, error) {
	config, err := configFromPath(filePath)

	if err != nil {
		return nil, err
	}

	c := &Client{
		config:  config,
		client:  &http.Client{Timeout: time.Duration(5 * time.Second)},
		BaseURL: adminURL(config),
	}

	return c, nil
}

func configFromPath(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := Config{}
	err = yaml.NewDecoder(file).Decode(&c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func adminURL(c *Config) string {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, c.Host)
}

// UpdateAllRecursively iterates through all services and updates config, deletes and recreates routes
func (c *Client) UpdateAllRecursively() error {
	for _, s := range c.config.Services {
		if err := c.UpdateService(s); err != nil {
			return err
		}
		if err := c.DeleteRoutes(s); err != nil {
			return err
		}
		if err := c.CreateRoutes(s); err != nil {
			return err
		}
	}
	return nil
}

// UpdateService updates a service's configuration in Kong
func (c *Client) UpdateService(s Service) error {
	url := fmt.Sprintf("%s/services/%s", c.BaseURL, s.Name)

	s.Routes = nil

	payload, err := json.Marshal(s)

	if err != nil {
		return err
	}

	res, err := c.httpRequest(http.MethodPut, url, payload, nil)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating service. Bad response from the API [%d]", res.StatusCode)
	}

	log.Printf("service [%s] created (%d)", s.Name, res.StatusCode)

	return nil
}

// DeleteRoutes iterates through the routes slice and deletes each one
func (c *Client) DeleteRoutes(s Service) error {
	routes, err := c.GetRoutes(s)
	if err != nil {
		return err
	}

	log.Printf("service [%s] - [%d] routes", s.Name, len(routes))

	for _, r := range routes {
		if err := c.DeleteRoute(r); err != nil {
			return err
		}
	}

	return nil
}

//
func (c *Client) DeleteRoute(r Route) error {
	url := fmt.Sprintf("%s/routes/%s", c.BaseURL, r.ID)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting route. Bad response response from the API [%d]", res.StatusCode)
	}

	log.Printf("route [%s] deleted (%d)", r.ID, res.StatusCode)

	return nil
}

func (c *Client) GetRoutes(s Service) ([]Route, error) {
	url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, s.Name)
	r := Routes{}
	res, err := c.httpRequest(http.MethodGet, url, nil, &r)
	if err != nil {
		return r.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("error fetching routes. Bad response response from the API [%d]", res.StatusCode)
	}

	return r.Data, nil
}

func (c *Client) CreateRoutes(s Service) error {
	url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, s.Name)

	for _, r := range s.Routes {
		payload, err := json.Marshal(r)
		if err != nil {
			return err
		}
		res, err := c.httpRequest(http.MethodPost, url, payload, nil)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusCreated {
			return fmt.Errorf("error creating routes. Bad response response from the API [%d]", res.StatusCode)
		}
	}

	log.Printf("routes created [%s]", s.Name)
	return nil
}

func (c *Client) CreatePlugin(s Service) error {
	url := fmt.Sprintf("%s/services/%s/plugins", c.BaseURL, s.Name)

	payload, err := json.Marshal(s.Plugin)
	if err != nil {
		return err
	}
	res, err := c.httpRequest(http.MethodPost, url, payload, nil)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating plugin. Bad response response from the API [%d]", res.StatusCode)
	}

	log.Printf("plugin created [%s]", s.Name)
	return nil
}

func (c *Client) httpRequest(method, url string, payload []byte, response interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set(contentType, applicationJSON)
	req.Header.Set("User-Agent", "kongfig")
	res, err := c.client.Do(req)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&response)

	// DEBUG
	// dump, _ := httputil.DumpResponse(res, true)
	// log.Printf("[http]: %q", dump)

	return res, err
}
