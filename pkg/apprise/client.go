package apprise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ClientConfig struct {
	Host string
	Key  string
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		Client: http.DefaultClient,
		config: config,
	}
}

type Client struct {
	config *ClientConfig
	*http.Client
}

func (c *Client) SendNotification(notification *Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/notify/%s", c.config.Host, c.config.Key), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
