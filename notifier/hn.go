package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

func CheckHN(q string) ([]Item, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "hn.algolia.com",
		Path:     "/api/v1/search",
		RawQuery: fmt.Sprintf("query=%s", url.QueryEscape(q)),
	}

	res, err := Client.Get(u.String())
	if err != nil {
		return nil, err
	}

	var algoliaRes algoliaResponse
	err = json.NewDecoder(res.Body).Decode(&algoliaRes)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, hit := range algoliaRes.Hits {
		if len(hit.ObjectID) == 0 {
			continue
		}
		isComment := hit.ParentID != 0
		if isComment {
			items = append(items, Item{
				ID:        hit.ObjectID,
				Title:     fmt.Sprintf("Comment on '%s'", hit.StoryTitle),
				Provider:  "Hacker News (comment)",
				URL:       fmt.Sprintf("https://news.ycombinator.com/item?id=%s", hit.ObjectID),
				Timestamp: hit.CreatedAt,
			})
		} else {
			items = append(items, Item{
				ID:        hit.ObjectID,
				Title:     hit.Title,
				Provider:  "Hacker News (story)",
				URL:       fmt.Sprintf("https://news.ycombinator.com/item?id=%d", hit.StoryID),
				Mention:   hit.URL,
				Timestamp: hit.CreatedAt,
			})
		}
	}

	return items, nil
}

type algoliaResponse struct {
	Exhaustive struct {
		NbHits bool `json:"nbHits"`
		Typo   bool `json:"typo"`
	} `json:"exhaustive"`
	ExhaustiveNbHits bool `json:"exhaustiveNbHits"`
	ExhaustiveTypo   bool `json:"exhaustiveTypo"`
	Hits             []struct {
		HighlightResult struct {
			Author struct {
				MatchLevel   string        `json:"matchLevel"`
				MatchedWords []interface{} `json:"matchedWords"`
				Value        string        `json:"value"`
			} `json:"author"`
			CommentText struct {
				FullyHighlighted bool     `json:"fullyHighlighted"`
				MatchLevel       string   `json:"matchLevel"`
				MatchedWords     []string `json:"matchedWords"`
				Value            string   `json:"value"`
			} `json:"comment_text"`
			StoryTitle struct {
				MatchLevel   string        `json:"matchLevel"`
				MatchedWords []interface{} `json:"matchedWords"`
				Value        string        `json:"value"`
			} `json:"story_title"`
			StoryURL struct {
				MatchLevel   string        `json:"matchLevel"`
				MatchedWords []interface{} `json:"matchedWords"`
				Value        string        `json:"value"`
			} `json:"story_url"`
		} `json:"_highlightResult"`
		Tags        []string  `json:"_tags"`
		Author      string    `json:"author"`
		CommentText string    `json:"comment_text,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		CreatedAtI  int       `json:"created_at_i"`
		ObjectID    string    `json:"objectID"`
		ParentID    int       `json:"parent_id,omitempty"`
		StoryID     int       `json:"story_id"`
		StoryTitle  string    `json:"story_title,omitempty"`
		StoryURL    string    `json:"story_url,omitempty"`
		UpdatedAt   time.Time `json:"updated_at"`
		NumComments int       `json:"num_comments,omitempty"`
		Points      int       `json:"points,omitempty"`
		Title       string    `json:"title,omitempty"`
		URL         string    `json:"url,omitempty"`
		Children    []int     `json:"children,omitempty"`
	} `json:"hits"`
	HitsPerPage         int    `json:"hitsPerPage"`
	NbHits              int    `json:"nbHits"`
	NbPages             int    `json:"nbPages"`
	Page                int    `json:"page"`
	Params              string `json:"params"`
	ProcessingTimeMS    int    `json:"processingTimeMS"`
	ProcessingTimingsMS struct {
		Request struct {
			RoundTrip int `json:"roundTrip"`
		} `json:"_request"`
		Fetch struct {
			Query    int `json:"query"`
			Scanning int `json:"scanning"`
			Total    int `json:"total"`
		} `json:"fetch"`
		Total int `json:"total"`
	} `json:"processingTimingsMS"`
	Query        string `json:"query"`
	ServerTimeMS int    `json:"serverTimeMS"`
}
