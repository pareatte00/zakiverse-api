package jikan

import (
	"context"
	"fmt"

	"github.com/zakiverse/zakiverse-api/util/http"
)

// --- Types ---

type Anime struct {
	MalID          int        `json:"mal_id"`
	URL            string     `json:"url"`
	Images         Images     `json:"images"`
	Trailer        Trailer    `json:"trailer"`
	Approved       bool       `json:"approved"`
	Titles         []Title    `json:"titles"`
	Title          string     `json:"title"`
	TitleEnglish   string     `json:"title_english"`
	TitleJapanese  string     `json:"title_japanese"`
	TitleSynonyms  []string   `json:"title_synonyms"`
	Type           string     `json:"type"`
	Source         string     `json:"source"`
	Episodes       *int       `json:"episodes"`
	Status         string     `json:"status"`
	Airing         bool       `json:"airing"`
	Aired          DateRange  `json:"aired"`
	Duration       string     `json:"duration"`
	Rating         string     `json:"rating"`
	Score          float64    `json:"score"`
	ScoredBy       int        `json:"scored_by"`
	Rank           int        `json:"rank"`
	Popularity     int        `json:"popularity"`
	Members        int        `json:"members"`
	Favorites      int        `json:"favorites"`
	Synopsis       string     `json:"synopsis"`
	Background     string     `json:"background"`
	Season         string     `json:"season"`
	Year           *int       `json:"year"`
	Broadcast      Broadcast  `json:"broadcast"`
	Producers      []MALEntry `json:"producers"`
	Licensors      []MALEntry `json:"licensors"`
	Studios        []MALEntry `json:"studios"`
	Genres         []MALEntry `json:"genres"`
	ExplicitGenres []MALEntry `json:"explicit_genres"`
	Themes         []MALEntry `json:"themes"`
	Demographics   []MALEntry `json:"demographics"`
}

type AnimeCharacter struct {
	Character   Character    `json:"character"`
	Role        string       `json:"role"`
	VoiceActors []VoiceActor `json:"voice_actors"`
}

type Character struct {
	MalID  int    `json:"mal_id"`
	URL    string `json:"url"`
	Images Images `json:"images"`
	Name   string `json:"name"`
}

type VoiceActor struct {
	Person   Person `json:"person"`
	Language string `json:"language"`
}

type AnimeStaffMember struct {
	Person    Person   `json:"person"`
	Positions []string `json:"positions"`
}

type AnimeEpisode struct {
	MalID         int    `json:"mal_id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	TitleJapanese string `json:"title_japanese"`
	TitleRomanji  string `json:"title_romanji"`
	Aired         string `json:"aired"`
	Score         *float64 `json:"score"`
	Filler        bool   `json:"filler"`
	Recap         bool   `json:"recap"`
	ForumURL      string `json:"forum_url"`
}

type AnimeVideo struct {
	Title   string  `json:"title"`
	Episode string  `json:"episode"`
	URL     string  `json:"url"`
	Images  Images  `json:"images"`
}

type AnimePromo struct {
	Title   string  `json:"title"`
	Trailer Trailer `json:"trailer"`
}

type AnimeStatistics struct {
	Watching    int `json:"watching"`
	Completed   int `json:"completed"`
	OnHold      int `json:"on_hold"`
	Dropped     int `json:"dropped"`
	PlanToWatch int `json:"plan_to_watch"`
	Total       int `json:"total"`
}

type AnimeThemes struct {
	Openings []string `json:"openings"`
	Endings  []string `json:"endings"`
}

type AnimeReview struct {
	MalID        int    `json:"mal_id"`
	URL          string `json:"url"`
	Type         string `json:"type"`
	Reactions    ReviewReactions `json:"reactions"`
	Date         string `json:"date"`
	Review       string `json:"review"`
	Score        int    `json:"score"`
	Tags         []string `json:"tags"`
	IsSpoiler    bool   `json:"is_spoiler"`
	IsPreliminary bool  `json:"is_preliminary"`
	EpisodesWatched int `json:"episodes_watched"`
}

type ReviewReactions struct {
	Overall   int `json:"overall"`
	Nice      int `json:"nice"`
	LoveIt    int `json:"love_it"`
	Funny     int `json:"funny"`
	Confusing int `json:"confusing"`
	Informative int `json:"informative"`
	WellWritten int `json:"well_written"`
	Creative  int `json:"creative"`
}

type ForumTopic struct {
	MalID         int    `json:"mal_id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	Date          string `json:"date"`
	AuthorUsername string `json:"author_username"`
	AuthorURL     string `json:"author_url"`
	Comments      int    `json:"comments"`
}

type NewsArticle struct {
	MalID        int    `json:"mal_id"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	AuthorUsername string `json:"author_username"`
	AuthorURL    string `json:"author_url"`
	ForumURL     string `json:"forum_url"`
	Images       Images `json:"images"`
	Comments     int    `json:"comments"`
	Excerpt      string `json:"excerpt"`
}

type Picture struct {
	JPG  ImageURLs `json:"jpg"`
	WebP ImageURLs `json:"webp"`
}

type UserUpdate struct {
	User         UserMeta `json:"user"`
	Score        *int     `json:"score"`
	Status       string   `json:"status"`
	EpisodesSeen *int     `json:"episodes_seen"`
	EpisodesTotal *int    `json:"episodes_total"`
	Date         string   `json:"date"`
}

type UserMeta struct {
	Username string `json:"username"`
	URL      string `json:"url"`
	Images   Images `json:"images"`
}

type Recommendation struct {
	Entry MALEntry `json:"entry"`
	Votes int      `json:"votes"`
}

// --- Endpoints ---

func (o *Client) GetAnime(ctx context.Context, id int) (*Anime, error) {
	var resp struct{ Data Anime `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetAnimeFull(ctx context.Context, id int) (*Anime, error) {
	var resp struct{ Data Anime `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/full", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetAnimeCharacters(ctx context.Context, id int) ([]AnimeCharacter, error) {
	var resp struct{ Data []AnimeCharacter `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/characters", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeStaff(ctx context.Context, id int) ([]AnimeStaffMember, error) {
	var resp struct{ Data []AnimeStaffMember `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/staff", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeEpisodes(ctx context.Context, id int, page int) ([]AnimeEpisode, *Pagination, error) {
	var resp struct {
		Data       []AnimeEpisode `json:"data"`
		Pagination Pagination     `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/episodes?page=%d", id, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}

func (o *Client) GetAnimeEpisode(ctx context.Context, id int, episode int) (*AnimeEpisode, error) {
	var resp struct{ Data AnimeEpisode `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/episodes/%d", id, episode)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetAnimeNews(ctx context.Context, id int, page int) ([]NewsArticle, *Pagination, error) {
	var resp struct {
		Data       []NewsArticle `json:"data"`
		Pagination Pagination    `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/news?page=%d", id, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}

func (o *Client) GetAnimeForum(ctx context.Context, id int) ([]ForumTopic, error) {
	var resp struct{ Data []ForumTopic `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/forum", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeVideos(ctx context.Context, id int) (struct {
	Promo    []AnimePromo `json:"promo"`
	Episodes []AnimeVideo `json:"episodes"`
}, error) {
	var resp struct {
		Data struct {
			Promo    []AnimePromo `json:"promo"`
			Episodes []AnimeVideo `json:"episodes"`
		} `json:"data"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/videos", id)})
	return resp.Data, err
}

func (o *Client) GetAnimePictures(ctx context.Context, id int) ([]Picture, error) {
	var resp struct{ Data []Picture `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/pictures", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeStatistics(ctx context.Context, id int) (*AnimeStatistics, error) {
	var resp struct{ Data AnimeStatistics `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/statistics", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetAnimeMoreInfo(ctx context.Context, id int) (string, error) {
	var resp struct{ Data struct{ Moreinfo string `json:"moreinfo"` } `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/moreinfo", id)})
	if err != nil { return "", err }
	return resp.Data.Moreinfo, nil
}

func (o *Client) GetAnimeRecommendations(ctx context.Context, id int) ([]Recommendation, error) {
	var resp struct{ Data []Recommendation `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/recommendations", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeUserUpdates(ctx context.Context, id int, page int) ([]UserUpdate, *Pagination, error) {
	var resp struct {
		Data       []UserUpdate `json:"data"`
		Pagination Pagination   `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/userupdates?page=%d", id, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}

func (o *Client) GetAnimeReviews(ctx context.Context, id int, page int) ([]AnimeReview, *Pagination, error) {
	var resp struct {
		Data       []AnimeReview `json:"data"`
		Pagination Pagination    `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/reviews?page=%d", id, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}

func (o *Client) GetAnimeRelations(ctx context.Context, id int) ([]Relation, error) {
	var resp struct{ Data []Relation `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/relations", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeThemes(ctx context.Context, id int) (*AnimeThemes, error) {
	var resp struct{ Data AnimeThemes `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/themes", id)})
	if err != nil { return nil, err }
	return &resp.Data, nil
}

func (o *Client) GetAnimeExternal(ctx context.Context, id int) ([]ExternalLink, error) {
	var resp struct{ Data []ExternalLink `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/external", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) GetAnimeStreaming(ctx context.Context, id int) ([]ExternalLink, error) {
	var resp struct{ Data []ExternalLink `json:"data"` }
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime/%d/streaming", id)})
	if err != nil { return nil, err }
	return resp.Data, nil
}

func (o *Client) SearchAnime(ctx context.Context, query string, page int) ([]Anime, *Pagination, error) {
	var resp struct {
		Data       []Anime    `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	_, err := o.client.Get(ctx, &resp, http.RequestParam{Path: fmt.Sprintf("/anime?q=%s&page=%d", query, page)})
	if err != nil { return nil, nil, err }
	return resp.Data, &resp.Pagination, nil
}
