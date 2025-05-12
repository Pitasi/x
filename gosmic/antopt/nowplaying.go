package antopt

import (
	"context"
	"encoding/json"
	"g2/antopt/lastfm"
	"iter"
	"log"
	"net/http"
	"sync"
	"time"
)

type nowPlayingSubs struct {
	mu   *sync.Mutex
	subs map[chan []byte]struct{}
	last []byte
}

func (s *nowPlayingSubs) subscribe(c chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[c] = struct{}{}
	c <- s.last
}

func (s *nowPlayingSubs) unsubscribe(c chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subs, c)
}

func (s *nowPlayingSubs) broadcast(next []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = next
	for c := range s.subs {
		c <- next
	}
}

func (ws Website) nowPlaying(c *lastfm.Client, mux *http.ServeMux) {
	subs := &nowPlayingSubs{
		mu:   &sync.Mutex{},
		subs: make(map[chan []byte]struct{}),
	}

	go func() {
		for track := range nowPlayingIter(c) {
			if track != nil {
				log.Println("now playing", track.Name, "-", track.Artist)
			}
			data, err := json.Marshal(track)
			if err != nil {
				log.Println("error marshalling now playing track:", err)
			} else {
				subs.broadcast(data)
			}
		}
	}()

	mux.HandleFunc("GET /now-playing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		w.WriteHeader(200)

		newTrack := make(chan []byte, 1)
		subs.subscribe(newTrack)
		defer func() {
			subs.unsubscribe(newTrack)
			close(newTrack)
		}()

		for {
			select {
			case <-r.Context().Done():
				return
			case data := <-newTrack:
				_, _ = w.Write([]byte("event: now-playing\n"))
				_, _ = w.Write([]byte("data: "))
				_, _ = w.Write(data)
				_, _ = w.Write([]byte("\n\n"))
				if w, ok := w.(http.Flusher); ok {
					w.Flush()
				}
			}
		}
	})
}

func nowPlayingIter(c *lastfm.Client) iter.Seq[*lastfm.Track] {
	var last *lastfm.Track
	return func(yield func(*lastfm.Track) bool) {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			track, err := c.NowPlaying(ctx)
			cancel()
			if err != nil {
				log.Println("fetching last.fm now playing:", err)
			} else if (last == nil && track != nil) || (last != nil && track == nil) || (last != nil && track != nil && *last != *track) {
				last = track
				if !yield(last) {
					return
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}
}
