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
	// Used in the creation of plugins for specific routes
	routeMap = make(map[string]string)
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

	if err := c.CreatePlugins(); err != nil {
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

// GetServices fetches all services from Kong
func (c *Client) GetServices() ([]Service, error) {
	url := fmt.Sprintf("%s/services", c.BaseURL)
	services := Services{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &services)

	if err != nil {
		return services.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return services.Data, fmt.Errorf("[HTTP %d] Error fetching Services. Bad response response from the API", res.StatusCode)
	}

	return services.Data, nil
}

// CreateRoutes iterates through all available routes and creates for the associated service
func (c *Client) CreateRoutes() error {
	for _, r := range c.config.Routes {
		url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, r.Service)

		payload, err := json.Marshal(r)

		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		if err != nil {
			return err
		}

		req.Header.Set(contentType, applicationJSON)
		req.Header.Set("User-Agent", "kongfig")

		res, err := c.client.Do(req)

		if err != nil {
			return err
		}

		defer res.Body.Close()

		route := Route{}
		json.NewDecoder(res.Body).Decode(&route)

		// Mapping route names to route ids
		// We do this so that we can create plugins for routes without having to
		// specific route id each time. It's easier to refer to routes via names
		routeMap[r.Name] = route.ID

		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("[HTTP %d] Error creating routes: Service not found", res.StatusCode)
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

// GetConsumers fetches all consumers from Kong
func (c *Client) GetConsumers() ([]Consumer, error) {
	url := fmt.Sprintf("%s/consumers", c.BaseURL)
	consumers := Consumers{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &consumers)

	if err != nil {
		return consumers.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return consumers.Data, fmt.Errorf("[HTTP %d] Error fetching routes. Bad response response from the API", res.StatusCode)
	}

	return consumers.Data, nil
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

// CreatePlugins creates global plugins, and plugins for services & routes
// Global plugins apply to all services and their routes
// Service plugins apply to all routes of a service
// Route plugins apply to only the specified route of a service
func (c *Client) CreatePlugins() error {
	for _, plugin := range c.config.Plugins {
		// Create global plugins
		if plugin.Target == "global" {
			url := fmt.Sprintf("%s/plugins", c.BaseURL)

			payload, err := json.Marshal(plugin)

			if err != nil {
				fmt.Println("Error marshalling payload: ", err)
				return err
			}

			res, err := c.httpRequest(http.MethodPost, url, payload, nil)

			if err != nil {
				fmt.Println("Error creating plugin: ", err)
				return err
			}

			if res.StatusCode == http.StatusNotFound {
				return fmt.Errorf("[HTTP %d] Global plugin already exists", res.StatusCode)
			}

			if res.StatusCode != http.StatusCreated {
				return fmt.Errorf("[HTTP %d] Error creating global plugin. Bad response from Kong API", res.StatusCode)
			}

			fmt.Printf("[HTTP %d] Global plugin created %s \n", res.StatusCode, plugin.Name)
		} else {
			// Creating plugins for specific services and routes
			// Create plugins for services:
			for _, service := range plugin.Services {
				url := fmt.Sprintf("%s/services/%s/plugins", c.BaseURL, service)

				payload, err := json.Marshal(plugin)

				if err != nil {
					fmt.Println("Error marshalling payload: ", err)
					return err
				}

				res, err := c.httpRequest(http.MethodPost, url, payload, nil)

				if err != nil {
					fmt.Println("Error creating plugin: ", err)
					return err
				}

				if res.StatusCode != http.StatusCreated {
					return fmt.Errorf("[HTTP %d] Error creating plugin for service %s. Bad response from Kong API", res.StatusCode, service)
				}

				fmt.Printf("[HTTP %d] Plugin created for service %s \n", res.StatusCode, service)
			}

			// Create plugins for routes
			for _, route := range plugin.Routes {
				routeID := routeMap[route]

				url := fmt.Sprintf("%s/routes/%s/plugins", c.BaseURL, routeID)
				payload, err := json.Marshal(plugin)

				if err != nil {
					fmt.Println("Error marshalling payload: ", err)
					return err
				}

				res, err := c.httpRequest(http.MethodPost, url, payload, nil)

				if err != nil {
					fmt.Println("Error creating plugin: ", err)
					return err
				}

				if res.StatusCode == http.StatusNotFound {
					return fmt.Errorf("[HTTP %d] Error creating plugin. Route not found %s", res.StatusCode, route)
				}

				if res.StatusCode != http.StatusCreated {
					return fmt.Errorf("[HTTP %d] Error creating plugin for route %s. Bad response from Kong API", res.StatusCode, route)
				}

				fmt.Printf("[HTTP %d] Plugin created for route %s \n", res.StatusCode, route)
			}
		}
	}

	return nil
}

// GetPlugins fetches all plugins from Kong
func (c *Client) GetPlugins() ([]Plugin, error) {
	url := fmt.Sprintf("%s/plugins", c.BaseURL)
	plugins := Plugins{}

	res, err := c.httpRequest(http.MethodGet, url, nil, &plugins)

	if err != nil {
		return plugins.Data, err
	}

	if res.StatusCode != http.StatusOK {
		return plugins.Data, fmt.Errorf("[HTTP %d] Error fetching Plugins. Bad response response from the API", res.StatusCode)
	}

	return plugins.Data, nil
}

// DeletePlugins iterates through all plugins and deletes each one
func (c *Client) DeletePlugins() error {
	plugins, err := c.GetPlugins()

	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if err := c.DeletePlugin(plugin); err != nil {
			return err
		}
	}

	return nil
}

// DeletePlugin deletes a plugin for a service based on route id
func (c *Client) DeletePlugin(plugin Plugin) error {
	url := fmt.Sprintf("%s/plugins/%s", c.BaseURL, plugin.ID)
	res, err := c.httpRequest(http.MethodDelete, url, nil, nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[HTTP %d] Error deleting plugin. Bad response response from the API", res.StatusCode)
	}

	fmt.Printf("[HTTP %d] Plugin [%s] deleted \n", res.StatusCode, plugin.Name)

	return nil
}
