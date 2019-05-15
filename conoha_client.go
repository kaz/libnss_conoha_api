package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	ConohaClient struct {
		http.Client
		Region   string
		Token    string
		Endpoint string
	}

	ServersResponse struct {
		Servers []Server
	}
	Server struct {
		Addresses map[string][]Address
		Metadata  Meta
	}
	Address struct {
		Version int
		Addr    string
	}
	Meta struct {
		Tag string `json:"instance_name_tag"`
	}
)

func NewClient(region, tenantId, username, password string) (*ConohaClient, error) {
	client := &ConohaClient{Region: region}

	data, err := json.Marshal(map[string]interface{}{
		"auth": map[string]interface{}{
			"passwordCredentials": map[string]string{
				"username": username,
				"password": password,
			},
			"tenantId": tenantId,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.Post("https://identity."+region+".conoha.io/v2.0/tokens", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	access := respData["access"].(map[string]interface{})
	token := access["token"].(map[string]interface{})
	serviceCatalog := access["serviceCatalog"].([]interface{})

	client.Token = token["id"].(string)

	for _, service := range serviceCatalog {
		svcMap := service.(map[string]interface{})
		if svcMap["type"].(string) == "compute" {
			client.Endpoint = svcMap["endpoints"].([]interface{})[0].(map[string]interface{})["publicURL"].(string)
			break
		}
	}
	if client.Endpoint == "" {
		return nil, fmt.Errorf("Failed to find compute endpoint")
	}

	return client, nil
}

func (cc *ConohaClient) get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", cc.Endpoint+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", cc.Token)
	resp, err := cc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

func (cc *ConohaClient) Servers() ([]Server, error) {
	var resp ServersResponse
	if err := cc.get("/servers/detail", &resp); err != nil {
		return nil, err
	}
	return resp.Servers, nil
}
