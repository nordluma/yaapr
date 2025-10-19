package allanime

import "context"

type episodes struct {
	Sub []string `json:"sub"`
	Dub []string `json:"dub"`
	Raw []string `json:"raw"`
}

type episodeResponse struct {
	Data struct {
		Show struct {
			AvailableEpisodesDetail episodes `json:"availableEpisodesDetail"`
		} `json:"show"`
	} `json:"data"`
}

type episodeVariables struct {
	ShowID string `json:"showId"`
}

func (c *AllanimeClient) GetEpisodes(
	ctx context.Context,
	showId string,
) (episodes, error) {
	query := `query ($showId: String!) { 
        show( _id: $showId ) {
            _id
            availableEpisodesDetail 
        }
    }`

	variables := episodeVariables{
		ShowID: showId,
	}

	var response episodeResponse
	if err := c.do(ctx, query, variables, &response); err != nil {
		return episodes{}, err
	}

	return response.Data.Show.AvailableEpisodesDetail, nil
}
