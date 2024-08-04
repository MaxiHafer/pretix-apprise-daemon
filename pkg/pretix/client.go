package pretix

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type transport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	if t.base == nil {
		t.base = http.DefaultTransport
	}

	return t.base.RoundTrip(req)
}

type ClientConfig struct {
	Token string
	Host  string
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		config: config,
		Client: &http.Client{
			Transport: &transport{
				headers: map[string]string{
					"Authorization": "Token " + config.Token,
				},
			},
		},
	}
}

type Client struct {
	*http.Client
	config *ClientConfig
}

func (c *Client) GetOrder(organizer, event, code string) (*Order, error) {
	resp, err := c.Client.Get(c.config.Host + fmt.Sprintf("/api/v1/organizers/%s/events/%s/orders/%s/", organizer, event, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	order := &Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}
