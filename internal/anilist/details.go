package anilist

import "context"

const getAnimeByIDQuery = `
query ($id: Int!) {
    Media(id: $id, type: ANIME) {
      id
      idMal
      title { romaji english native }
      status(version: 2)
      genres
      episodes
      coverImage { large }
      averageScore
      meanScore
    }
}`

type AnimeDetails struct {
	ID    int `json:"id"`
	IDMal int `json:"idMal"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	Status     string   `json:"status"` // TODO: use `MediaStatus`
	Genres     []string `json:"genres"`
	Episodes   int      `json:"episodes"`
	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`
	AverageScore int `json:"averageScore"`
	MeanScore    int `json:"meanScore"`
}

type getAnimeByIdResponse struct {
	Data struct {
		Media struct {
			AnimeDetails
		} `json:"media"`
	} `json:"data"`
}

func (c *AnilistClient) GetAnimeById(
	ctx context.Context,
	id int,
) (AnimeDetails, error) {
	var res getAnimeByIdResponse
	err := c.do(ctx, getAnimeByIDQuery, map[string]any{"id": id}, &res)
	if err != nil {
		return AnimeDetails{}, err
	}

	return res.Data.Media.AnimeDetails, nil
}
