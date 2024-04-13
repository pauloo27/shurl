package models

type Link struct {
	Slug        string `redis:"-" json:"slug"`
	Domain      string `redis:"-" json:"domain"`
	OriginalURL string `redis:"url" json:"original_url"`
	TTL         int    `redis:"-" json:"ttl"`
}
