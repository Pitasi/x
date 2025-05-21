package antopt

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"sync"
	"time"

	"anto.pt/x/gosmic/antopt/lastfm"
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
				logger.Info("now playing", "name", track.Name, "artist", track.Artist)
			}
			data, err := json.Marshal(track)
			if err != nil {
				logger.Error("marshalling now playing track", "err", err)
			} else {
				subs.broadcast(data)
			}
		}
	}()

	keepalive := 30 * time.Second
	mux.HandleFunc("GET /now-playing", func(w http.ResponseWriter, r *http.Request) {
		ctrl := http.NewResponseController(w)

		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		w.WriteHeader(200)

		newTrack := make(chan []byte, 1)
		subs.subscribe(newTrack)
		defer func() {
			subs.unsubscribe(newTrack)
			close(newTrack)
		}()

		keepaliveTimer := time.NewTicker(keepalive)
		for {
			select {
			case <-r.Context().Done():
				return
			case <-keepaliveTimer.C:
				_, _ = w.Write([]byte("data: ping"))
				_, _ = w.Write([]byte("\n\n"))
				_ = ctrl.Flush()
			case data := <-newTrack:
				_, _ = w.Write([]byte("event: now-playing\n"))
				_, _ = w.Write([]byte("data: "))
				_, _ = w.Write(data)
				_, _ = w.Write([]byte("\n\n"))
				_ = ctrl.Flush()
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
				logger.Error("fetching last.fm now playing", "err", err)
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
