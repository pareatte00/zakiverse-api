package discord

import (
	"context"

	"github.com/zakiverse/zakiverse-api/util/http"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (o *Client) GetUser(ctx context.Context, accessToken string) (User, error) {
	var resp User

	_, err := o.client.Get(ctx, &resp, http.RequestParam{
		Path: "/users/@me",
		Header: map[string]string{
			"Authorization": "Bearer " + accessToken,
		},
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
