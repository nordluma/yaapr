package anilist

import "context"

const searchQuery = `
query ($search: String) {
  Page(page: 1, perPage: 10) {
    media(search: $search, type: ANIME) {
      id
      title { romaji english native }
      episodes
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
	Episodes   int `json:"episodes"`
	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`
}

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
