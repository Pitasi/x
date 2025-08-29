package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func main() {
	for _, arg := range os.Args[1:] {
		process(arg)
	}
}

func process(path string) {
	log.Println("processing", path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Println("ignoring unsupported file", path)
		return
	}

	u, _ := url.Parse("https://lrclib.net/api/get")
	q := u.Query()
	q.Set("artist_name", strings.ReplaceAll(m.Artist(), " ", "+"))
	q.Set("track_name", strings.ReplaceAll(m.Title(), " ", "+"))
	q.Set("album_name", strings.ReplaceAll(m.Album(), " ", "+"))
	q.Set("duration", fmt.Sprintf("%0f", m.Duration().Seconds()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("user-agent", "lyr v0 (https://github.com/Pitasi/x/tree/main/lyr)")

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("fetch", u.String(), ": status code", res.Status)
	}

	type response struct {
		PlainLyrics  string `json:"plainLyrics"`
		SyncedLyrics string `json:"syncedLyrics"`
	}
	var resPayload response
	if err := json.NewDecoder(res.Body).Decode(&resPayload); err != nil {
		log.Fatal(err)
	}

	if resPayload.PlainLyrics == "" && resPayload.SyncedLyrics == "" {
		log.Fatal("no lyrics")
	}

	ext := filepath.Ext(path)

	if resPayload.SyncedLyrics != "" {
		lrc := strings.Replace(path, ext, ".lrc", 1)
		if err := os.WriteFile(lrc, []byte(resPayload.SyncedLyrics), 0666); err != nil {
			log.Fatal(err)
		}
		log.Println("written", lrc)
	} else {
		txt := strings.Replace(path, ext, ".txt", 1)
		if err := os.WriteFile(txt, []byte(resPayload.PlainLyrics), 0666); err != nil {
			log.Fatal(err)
		}
		log.Println("written", txt)
	}
}
