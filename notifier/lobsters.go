package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"
)

func CheckLobsters(domain string) ([]Item, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "lobste.rs",
		Path:   "/domains/" + domain + ".json",
	}

	res, err := Client.Get(u.String())
	if err != nil {
		return nil, err
	}

	var lobstersRes []lobstersResponse
	err = json.NewDecoder(res.Body).Decode(&lobstersRes)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, res := range lobstersRes {
		items = append(items, Item{
			ID:        res.ShortID,
			Title:     res.Title,
			Provider:  "Lobsters",
			URL:       res.ShortIDURL,
			Mention:   res.URL,
			Timestamp: parseTimestamp(res.CreatedAt),
		})
	}

	return items, nil
}

func parseTimestamp(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000-07:00", s)
	if err != nil {
		log.Println("error parsing lobsters timestamp:", s)
	}
	return t
}

type lobstersResponse struct {
	ShortID          string   `json:"short_id"`
	ShortIDURL       string   `json:"short_id_url"`
	CreatedAt        string   `json:"created_at"`
	Title            string   `json:"title"`
	URL              string   `json:"url"`
	Score            int      `json:"score"`
	Flags            int      `json:"flags"`
	CommentCount     int      `json:"comment_count"`
	Description      string   `json:"description"`
	DescriptionPlain string   `json:"description_plain"`
	CommentsURL      string   `json:"comments_url"`
	SubmitterUser    string   `json:"submitter_user"`
	UserIsAuthor     bool     `json:"user_is_author"`
	Tags             []string `json:"tags"`
}
