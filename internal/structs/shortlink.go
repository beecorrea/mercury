package structs

import "time"

type Shortlink struct {
	Key           string    `json:"key"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Summary       string    `json:"summary"`
	ScrapeAttempts int      `json:"scrape_attempts"`
	CreatedAt     time.Time `json:"created_at"`
}

