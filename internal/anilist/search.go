package anilist

import "context"

const searchQuery = `
query ($search: String) {
  Page(page: 1, perPage: 10) {
    media(search: $search, type: ANIME) {
      id
      title { romaji english native }
      coverImage { large }
    }
  }
}`

type Anime struct {
	ID    int `json:"id"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`
}

func (a Anime) FilterValue() string { return "" }

type searchResponse struct {
	Data struct {
		Page struct {
			Media []Anime `json:"media"`
		} `json:"page"`
	} `json:"data"`
}

func (c *Client) SearchAnime(
	ctx context.Context,
	name string,
) ([]Anime, error) {
	var res searchResponse
	err := c.do(ctx, searchQuery, map[string]any{"search": name}, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.Page.Media, nil
}
