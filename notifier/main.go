package main

import (
	"fmt"
	"log"
	"time"
)

type Item struct {
	ID        string
	Title     string
	Provider  string
	Mention   string
	URL       string
	Timestamp time.Time
}

func main() {
	for {
		log.Println("running checks")

		hn, err := CheckHN("anto.pt")
		if err != nil {
			SendError(fmt.Errorf("error checking hn: %w", err))
		}

		lobsters, err := CheckLobsters("anto.pt")
		if err != nil {
			SendError(fmt.Errorf("error checking lobsters: %w", err))
		}

		items := dedupe(append(hn, lobsters...))

		if len(items) == 0 {
			log.Println("no new items")
		} else {
			log.Println("sending", len(items), "new items")
		}

		SendItems(items)

		time.Sleep(time.Hour * 6)
	}
}

var seen = make(map[string]struct{})

func dedupe(items []Item) []Item {
	var deduped []Item
	for _, item := range items {
		if _, s := seen[item.ID]; s {
			continue
		}
		seen[item.ID] = struct{}{}
		deduped = append(deduped, item)
	}
	return deduped
}
