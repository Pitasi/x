package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var User = "zaphodias"

type Client struct {
	APIKey string
}

type Track struct {
	Name   string
	Artist string
	Image  template.URL
}

func (c *Client) NowPlaying(ctx context.Context) (*Track, error) {
	res, err := c.recentTracks(ctx)
	if err != nil {
		return nil, fmt.Errorf("api error: %w", err)
	}

	if len(res.RecentTracks.Track) == 0 {
		return nil, nil
	}

	firstTrack := res.RecentTracks.Track[0]
	uts, err := strconv.ParseInt(firstTrack.Date.UTS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing unix timestamp string: %w", err)
	}
	ts := time.Unix(uts, 0)

	if time.Now().Add(-10 * time.Minute).After(ts) {
		// track was listened more than 10 minutes ago, is not currently playing
		return nil, nil
	}

	var image string
	for _, i := range firstTrack.Image {
		if i.Size == "medium" {
			image = i.Text
		}
	}
	if image == "" && len(firstTrack.Image) > 0 {
		image = firstTrack.Image[0].Text
	}

	return &Track{
		Name:   firstTrack.Name,
		Artist: firstTrack.Artist.Text,
		Image:  template.URL(image),
	}, nil
}

type ArtistResponse struct {
	Text string `json:"#text"`
}

type ImageResponse struct {
	Size string `json:"size"`
	Text string `json:"#text"`
}

type DateResponse struct {
	UTS string `json:"uts"`
}

type TrackResponse struct {
	Artist ArtistResponse  `json:"artist"`
	Image  []ImageResponse `json:"image"`
	Name   string          `json:"name"`
	Date   DateResponse    `json:"date"`
}

type RecentTracksResponse struct {
	RecentTracks struct {
		Track []TrackResponse `json:"track"`
	} `json:"recenttracks"`
}

func (c *Client) recentTracks(ctx context.Context) (RecentTracksResponse, error) {
	u, err := url.Parse("https://ws.audioscrobbler.com/2.0")
	if err != nil {
		return RecentTracksResponse{}, fmt.Errorf("url parse: %w", err)
	}

	q := u.Query()
	q.Set("method", "user.getrecenttracks")
	q.Set("format", "json")
	q.Set("user", User)
	q.Set("limit", "5")
	q.Set("api_key", c.APIKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return RecentTracksResponse{}, fmt.Errorf("preparing request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RecentTracksResponse{}, fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		response, err := io.ReadAll(res.Body)
		if err != nil {
			return RecentTracksResponse{}, fmt.Errorf("api status: %s, couldn't read body: %w", res.Status, err)
		}
		return RecentTracksResponse{}, fmt.Errorf("api status: %s, body: %s", res.Status, string(response))
	}

	var body RecentTracksResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return RecentTracksResponse{}, fmt.Errorf("parsing json response: %w", err)
	}

	return body, nil
}
