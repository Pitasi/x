// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from encoding/xml/read_test.go.
package antopt

import (
	"encoding/xml"
	"errors"
	"iter"
	"time"

	"anto.pt/x/gosmic/antopt/articles"
)

type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Link    []Link   `xml:"link"`
	Updated TimeStr  `xml:"updated"`
	Author  *Person  `xml:"author"`
	Entry   []*Entry `xml:"entry"`
}

type Entry struct {
	Title     string  `xml:"title"`
	ID        string  `xml:"id"`
	Link      []Link  `xml:"link"`
	Published TimeStr `xml:"published"`
	Updated   TimeStr `xml:"updated"`
	Author    *Person `xml:"author"`
	Summary   *Text   `xml:"summary"`
	Content   *Text   `xml:"content"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri"`
	Email    string `xml:"email"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

func buildArticlesAtomFeed(a articles.Articles) ([]byte, error) {
	next, stop := iter.Pull2(a.Published())
	_, latest, ok := next()
	if !ok {
		return nil, errors.New("no articles")
	}
	stop()

	feed := Feed{
		Title: "Antonio Pitasi",
		ID:    "https://anto.pt/",
		Link: []Link{
			{Rel: "self", Href: "https://anto.pt/articles/feed.atom"},
		},
		Author: &Person{
			Name:  "Antonio Pitasi",
			URI:   "https://anto.pt/",
			Email: "antonio@pitasi.dev",
		},
		Updated: Time(latest.Date),
		Entry:   articlesToAtomEntries(a),
	}
	return xml.Marshal(feed)
}

func articlesToAtomEntries(a articles.Articles) []*Entry {
	var (
		max     = 10
		entries = make([]*Entry, 0, max)
	)

	for _, a := range limit2(a.Published(), max) {
		href := "https://anto.pt/articles/" + a.Slug
		entry := &Entry{
			Title: a.Title,
			ID:    href,
			Link: []Link{
				{Rel: "alternate", Href: href},
			},
			Published: Time(a.Date),
			Updated:   Time(a.Date),
			Summary: &Text{
				Type: "text",
				Body: a.Description,
			},
			Content: &Text{
				Type: "html",
				Body: string(a.Content),
			},
		}

		entries = append(entries, entry)
	}

	return entries
}

func limit2[K, V any](it iter.Seq2[K, V], limit int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var count int
		for k, v := range it {
			if !yield(k, v) {
				return
			}
			count++
			if count == limit {
				return
			}
		}
	}
}
