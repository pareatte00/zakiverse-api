package discord

import (
	"context"
	"net/url"

	"github.com/zakiverse/zakiverse-api/util/http"
)

type ExchangeCodeParam struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
	Code         string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (o *Client) ExchangeCode(ctx context.Context, param ExchangeCodeParam) (TokenResponse, error) {
	var resp TokenResponse

	_, err := o.client.Post(ctx, &resp, http.RequestParam{
		Path:     "/oauth2/token",
		BodyType: http.BodyForm,
		Body: url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {param.Code},
			"redirect_uri":  {param.RedirectURI},
			"client_id":     {param.ClientId},
			"client_secret": {param.ClientSecret},
		},
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
