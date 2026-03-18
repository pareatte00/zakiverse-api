package jikan

import (
	"context"
	"fmt"

	"github.com/zakiverse/zakiverse-api/util/http"
)

// --- Types ---

type CharacterFull struct {
	MalID    int      `json:"mal_id"`
	URL      string   `json:"url"`
	Images   Images   `json:"images"`
	Name     string   `json:"name"`
	NameKanji string  `json:"name_kanji"`
	Nicknames []string `json:"nicknames"`
	Favorites int     `json:"favorites"`
	About    string   `json:"about"`
}

type CharacterAnime struct {
	Role  string `json:"role"`
	Anime Anime  `json:"anime"`
}

type CharacterManga struct {
	Role  string `json:"role"`
	Manga struct {
		MalID  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images Images `json:"images"`
		Title  string `json:"title"`
	} `json:"manga"`
}

type CharacterVoiceActor struct {
	Language string `json:"language"`
	Person   Person `json:"person"`
}

// --- Endpoints ---

func (o *Client) GetCharacter(ctx context.Context, id int) (*CharacterFull, error) {
	var resp struct{ Data CharacterFull `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetCharacterFull(ctx context.Context, id int) (*CharacterFull, error) {
	var resp struct{ Data CharacterFull `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d/full", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetCharacterAnime(ctx context.Context, id int) ([]CharacterAnime, error) {
	var resp struct{ Data []CharacterAnime `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d/anime", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetCharacterManga(ctx context.Context, id int) ([]CharacterManga, error) {
	var resp struct{ Data []CharacterManga `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d/manga", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetCharacterVoices(ctx context.Context, id int) ([]CharacterVoiceActor, error) {
	var resp struct{ Data []CharacterVoiceActor `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d/voices", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetCharacterPictures(ctx context.Context, id int) ([]Picture, error) {
	var resp struct{ Data []Picture `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters/%d/pictures", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) SearchCharacters(ctx context.Context, query string, page int) ([]CharacterFull, *Pagination, error) {
	var resp struct {
		Data       []CharacterFull `json:"data"`
		Pagination Pagination      `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/characters?q=%s&page=%d", query, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}
