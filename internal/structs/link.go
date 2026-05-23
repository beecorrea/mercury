package structs

import "time"

type Link struct {
	Key       string    `json:"key"`
	URL       string    `json:"url"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
}
