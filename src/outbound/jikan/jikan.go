package jikan

import (
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/util/http"
)

type Client struct {
	client *http.Client
}

func New(conf config.ConfigConstant) *Client {
	return &Client{
		client: http.MustNew(conf.Outbound.Jikan.BaseUrl),
	}
}

func (o *Client) Disconnect() {
	o.client.Disconnect()
}
