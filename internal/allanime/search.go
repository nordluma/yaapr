package allanime

import (
	"context"
)

type anime struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	EnglishName       string            `json:"englishName"`
	AvailableEpisodes availableEpisodes `json:"availableEpisodes"`
}

type availableEpisodes struct {
	Sub int `json:"sub"`
	Dub int `json:"dub"`
	Raw int `json:"raw"`
}

type searchResponse struct {
	Data struct {
		Shows struct {
			Edges []anime `json:"edges"`
		} `json:"shows"`
	} `json:"data"`
}

type searchVariables struct {
	Search          searchOpts `json:"search"`
	Limit           int        `json:"limit"`
	Page            int        `json:"page"`
	TranslationType string     `json:"translationType"`
	CountryOrigin   string     `json:"countryOrigin"`
}

type searchOpts struct {
	AllowAdult   bool   `json:"allowAdult"`
	AllowUnknown bool   `json:"allowUnknown"`
	Query        string `json:"query"`
}

func (c *AllanimeClient) SearchAnime(
	ctx context.Context,
	mode, query string,
) ([]anime, error) {
	searchQuery := `query($search: SearchInput, $limit: Int, $page: Int, $translationType: VaildTranslationTypeEnumType, $countryOrigin: VaildCountryOriginEnumType) {
		shows(search: $search, limit: $limit, page: $page, translationType: $translationType, countryOrigin: $countryOrigin) {
			edges {
				_id
				name
				englishName
				availableEpisodes
				__typename
			}
		}
	}`

	variables := searchVariables{
		Search: searchOpts{
			AllowAdult:   false,
			AllowUnknown: false,
			Query:        query,
		},
		Limit:           40,
		Page:            1,
		TranslationType: mode,
		CountryOrigin:   "ALL",
	}

	var searchRes searchResponse
	if err := c.do(ctx, searchQuery, variables, &searchRes); err != nil {
		return nil, err
	}

	return searchRes.Data.Shows.Edges, nil
}
