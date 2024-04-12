package models

type Link struct {
	Slug   string `redis:"-" json:"slug"`
	Domain string `redis:"-" json:"domain"`
	URL    string `redis:"url" json:"url"`
	TTL    int    `redis:"-" json:"ttl"`
}
