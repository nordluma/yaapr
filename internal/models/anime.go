package models

type Anime struct {
	ID    string     `json:"id"`
	Title AnimeTitle `json:"title"`
	// MalId       int        `json:"idMal"`
	// Description string     `json:"description"`
}

type AnimeTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	// Native  string `json:"native"`
}
