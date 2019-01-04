package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// httpRequest is an utility method for executing HTTP requests
func (c *Client) httpRequest(method, url string, payload []byte, response interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set(contentType, applicationJSON)
	req.Header.Set("User-Agent", "kongfig")
	res, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&response)

	return res, err
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

// configFromPath parses the YAML file specified in the path param
func configFromPath(path string) (*Config, error) {
	configData, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	configData = []byte(os.ExpandEnv(string(configData)))

	c := Config{}

	if err := yaml.Unmarshal(configData, &c); err != nil {
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

// ApplyConfig iterates through all services and updates config, deletes and recreates routes
func (c *Client) ApplyConfig() error {

	if err := c.DeleteConsumers(); err != nil {
		return err
	}

	if err := c.DeleteRoutes(); err != nil {
		return err
	}

	if err := c.DeleteServices(); err != nil {
		return err
	}

	if err := c.DeletePlugins(); err != nil {
		return err
	}

	for _, s := range c.config.Services {
		if err := c.UpdateService(s); err != nil {
			return err
		}

	}

	if err := c.CreateRoutes(); err != nil {
		return err
	}

	return nil
}

// UpdateService updates an existing service or creates a new one if it doesn't exist
// Makes a HTTP PUT to the KONG ADMIN API
func (c *Client) UpdateService(s Service) error {
	url := fmt.Sprintf("%s/services/%s", c.BaseURL, s.Name)

	payload, err := json.Marshal(s)

	if err != nil {
		return err
	}

	res, err := c.httpRequest(http.MethodPut, url, payload, nil)

	if err != nil {
		fmt.Printf("error: %s \n", err)
		fmt.Printf("Error updating service: %s \n", s.Name)

		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("[HTTP %d] Error updating service. Bad response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Successfully created/updated service: %s \n", http.StatusOK, s.Name)

	return nil
}

// CreateRoutes iterates through all available routes and creates for the associated service
func (c *Client) CreateRoutes() error {
	for _, r := range c.config.Routes {
		url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, r.Service)

		payload, err := json.Marshal(r)

		if err != nil {
			return err
		}

		res, err := c.httpRequest(http.MethodPost, url, payload, nil)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusCreated {
			return fmt.Errorf("[HTTP %d] Error creating routes. Bad response from Kong API", res.StatusCode)
		}

		fmt.Printf("[HTTP %d] Route created for service %s \n", res.StatusCode, r.Service)
	}

	return nil
}

// GetRoutes fetches all routes from Kong
func (c *Client) GetRoutes() ([]Route, error) {
	url := fmt.Sprintf("%s/routes", c.BaseURL)
	r := Routes{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &r)

	if err != nil {
		return r.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("[HTTP %d] Error fetching routes. Bad response response from the API", res.StatusCode)
	}

	return r.Data, nil
}

// DeleteRoutes iterates through all routes and deletes each one
func (c *Client) DeleteRoutes() error {
	routes, err := c.GetRoutes()

	if err != nil {
		return err
	}

	for _, r := range routes {
		if err := c.DeleteRoute(r); err != nil {
			return err
		}
	}

	return nil
}

// DeleteRoute deletes a route for a service based on route id
func (c *Client) DeleteRoute(r Route) error {
	url := fmt.Sprintf("%s/routes/%s", c.BaseURL, r.ID)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[HTTP %d] Error deleting route. Bad response response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Route [%s] deleted \n", res.StatusCode, r.ID)

	return nil
}

// GetCredentials fetches all consumers from Kong
func (c *Client) GetConsumers() ([]Consumer, error) {
	url := fmt.Sprintf("%s/consumers", c.BaseURL)
	r := Consumers{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &r)

	if err != nil {
		fmt.Println("Here")
		return r.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("[HTTP %d] Error fetching consumers. Bad response response from the API", res.StatusCode)
	}

	return r.Data, nil
}

// DeleteConsumers iterates through all routes and deletes each one
func (c *Client) DeleteConsumers() error {
	consumers, err := c.GetConsumers()

	if err != nil {
		return err
	}

	for _, r := range consumers {
		if err := c.DeleteConsumer(r); err != nil {
			return err
		}
	}

	return nil
}

// DeleteConsumer deletes a consumer for a service based on route id
func (c *Client) DeleteConsumer(r Consumer) error {
	url := fmt.Sprintf("%s/consumers/%s", c.BaseURL, r.Username)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[HTTP %d] Error deleting consumer. Bad response response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Consumer [%s] deleted \n", res.StatusCode, r.Username)

	return nil
}

// GetPlugins fetches all plugins from Kong
func (c *Client) GetPlugins() ([]Plugin, error) {
	url := fmt.Sprintf("%s/plugins", c.BaseURL)
	r := Plugins{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &r)

	if err != nil {
		fmt.Println("Here")
		return r.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("[HTTP %d] Error fetching Plugins. Bad response response from the API", res.StatusCode)
	}

	return r.Data, nil
}

// DeletePlugins iterates through all plugins and deletes each one
func (c *Client) DeletePlugins() error {
	plugins, err := c.GetPlugins()

	if err != nil {
		return err
	}

	for _, r := range plugins {
		if err := c.DeletePlugin(r); err != nil {
			return err
		}
	}

	return nil
}

// DeletePlugin deletes a plugin for a service based on route id
func (c *Client) DeletePlugin(r Plugin) error {
	url := fmt.Sprintf("%s/plugins/%s", c.BaseURL, r.Name)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[HTTP %d] Error deleting plugin. Bad response response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Plugin [%s] deleted \n", res.StatusCode, r.Name)

	return nil
}

// GetServices fetches all services from Kong
func (c *Client) GetServices() ([]Service, error) {
	url := fmt.Sprintf("%s/services", c.BaseURL)
	r := Services{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &r)

	if err != nil {
		return r.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("[HTTP %d] Error fetching Services. Bad response response from the API", res.StatusCode)
	}

	return r.Data, nil
}

// DeleteServices iterates through all services and deletes each one
func (c *Client) DeleteServices() error {
	services, err := c.GetServices()

	if err != nil {
		return err
	}

	for _, r := range services {
		if err := c.DeleteService(r); err != nil {
			return err
		}
	}

	return nil
}

// DeleteService deletes a service for a service based on route id
func (c *Client) DeleteService(r Service) error {
	url := fmt.Sprintf("%s/services/%s", c.BaseURL, r.Name)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[HTTP %d] Error deleting service. Bad response response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Service [%s] deleted \n", res.StatusCode, r.Name)

	return nil
}

// Commenting out for now. A future version of the Plugins feature will replace this:
// func (c *Client) CreatePlugin(s Service) error {
// 	url := fmt.Sprintf("%s/services/%s/plugins", c.BaseURL, s.Name)

// 	payload, err := json.Marshal(s.Plugin)
// 	if err != nil {
// 		return err
// 	}
// 	res, err := c.httpRequest(http.MethodPost, url, payload, nil)
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("error creating plugin. Bad response response from the API [%d]", res.StatusCode)
// 	}

// 	log.Printf("plugin created [%s]", s.Name)
// 	return nil
// }
