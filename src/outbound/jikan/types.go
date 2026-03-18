package jikan

// Common types shared across endpoints

type Pagination struct {
	LastVisiblePage int  `json:"last_visible_page"`
	HasNextPage     bool `json:"has_next_page"`
	CurrentPage     int  `json:"current_page"`
	Items           struct {
		Count   int `json:"count"`
		Total   int `json:"total"`
		PerPage int `json:"per_page"`
	} `json:"items"`
}

type ImageURLs struct {
	ImageURL      string `json:"image_url"`
	SmallImageURL string `json:"small_image_url"`
	LargeImageURL string `json:"large_image_url"`
}

type Images struct {
	JPG  ImageURLs `json:"jpg"`
	WebP ImageURLs `json:"webp"`
}

type MALEntry struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Title struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type DateRange struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Prop   DateProp `json:"prop"`
	String string `json:"string"`
}

type DateProp struct {
	From DateParts `json:"from"`
	To   DateParts `json:"to"`
}

type DateParts struct {
	Day   *int `json:"day"`
	Month *int `json:"month"`
	Year  *int `json:"year"`
}

type Trailer struct {
	YoutubeID string        `json:"youtube_id"`
	URL       string        `json:"url"`
	EmbedURL  string        `json:"embed_url"`
	Images    TrailerImages `json:"images"`
}

type TrailerImages struct {
	ImageURL        string `json:"image_url"`
	SmallImageURL   string `json:"small_image_url"`
	MediumImageURL  string `json:"medium_image_url"`
	LargeImageURL   string `json:"large_image_url"`
	MaximumImageURL string `json:"maximum_image_url"`
}

type Broadcast struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
	String   string `json:"string"`
}

type Relation struct {
	Relation string     `json:"relation"`
	Entry    []MALEntry `json:"entry"`
}

type ExternalLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Person struct {
	MalID  int    `json:"mal_id"`
	URL    string `json:"url"`
	Images Images `json:"images"`
	Name   string `json:"name"`
}
