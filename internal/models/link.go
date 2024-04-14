package models

type Link struct {
	Slug        string `json:"slug"`
	Domain      string `json:"domain"`
	OriginalURL string `json:"original_url"`
	URL         string `json:"url"`
	TTL         int    `json:"ttl"`
}
