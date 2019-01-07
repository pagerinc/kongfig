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

var (
	// Keeps track of route names to route IDs
	routeMap = make(map[string]string)
)

// type response struct {
// 	ID string `json:"id"`
// }

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

	// defer res.Body.Close()
	// json.NewDecoder(res.Body).Decode(&response)

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
	for _, s := range c.config.Services {
		if err := c.UpdateService(s); err != nil {
			return err
		}

		if err := c.DeleteRoutes(s); err != nil {
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
		fmt.Println("PAYLOAD: ", string(payload))

		if err != nil {
			return err
		}

		res, err := c.httpRequest(http.MethodPost, url, payload, nil)

		var response interface{}

		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&response)

		fmt.Println("RESPONSE BODY: ", response)

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

// GetRoutes fetches all routes from Kong for the specified service
func (c *Client) GetRoutes(s Service) ([]Route, error) {
	url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, s.Name)
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
func (c *Client) DeleteRoutes(s Service) error {
	routes, err := c.GetRoutes(s)

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

// CreatePlugins creates plugins for associated services
func (c *Client) CreatePlugins(service Service) error {
	for _, plugin := range c.config.Plugins {

		for _, service := range plugin.Services {
			url := fmt.Sprintf("%s/services/%s/plugins", c.BaseURL, service)

			payload, err := json.Marshal(plugin)

			if err != nil {
				return err
			}

			res, err := c.httpRequest(http.MethodPost, url, payload, nil)

			if err != nil {
				return err
			}

			if res.StatusCode != http.StatusCreated {
				return fmt.Errorf("[HTTP %d] Error creating plugin. Bad response from Kong API", res.StatusCode)
			}

			fmt.Printf("[HTTP %d] Plugin created for service %s \n", res.StatusCode, service)
		}

		// for _, route := range plugin.Routes {
		// 	url := fmt.Sprintf("%s/services/%s/plugins", c.BaseURL, service)

		// 	payload, err := json.Marshal(plugin)

		// 	if err != nil {
		// 		return err
		// 	}

		// 	res, err := c.httpRequest(http.MethodPost, url, payload, nil)

		// 	if err != nil {
		// 		return err
		// 	}

		// 	if res.StatusCode != http.StatusCreated {
		// 		return fmt.Errorf("[HTTP %d] Error creating plugin. Bad response from Kong API", res.StatusCode)
		// 	}

		// 	fmt.Printf("[HTTP %d] Plugin created for service %s \n", res.StatusCode, service)
		// }
	}

	return nil
}

// func (c *Client) DeletePlugin(service Service) error {

// }

// func (c *Client) DeletePlugins(service Service) error {

// }
