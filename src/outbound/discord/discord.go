package discord

import (
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/util/http"
)

type Client struct {
	client *http.Client
}

func New(conf config.ConfigConstant) *Client {
	return &Client{
		client: http.MustNew(conf.Outbound.Discord.BaseUrl),
	}
}

func (o *Client) Disconnect() {
	o.client.Disconnect()
}
