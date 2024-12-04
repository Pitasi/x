package main

import "html/template"

type Site struct {
	Nav            []NavItem
	StaticFileHash map[string]string
	Colors         []template.CSS
}

type NavItem struct {
	Title string
	URL   string
}

type ArticlesByYear struct {
	Year  string
	Posts []ArticleLink
}

type ArticleLink struct {
	Title       string
	Description string
	Date        string
	URL         string
}

var site Site

func init() {
	site = Site{
		Nav: []NavItem{
			{
				Title: "Articles",
				URL:   "/",
			},
			{
				Title: "Uses",
				URL:   "/uses",
			},
			{
				Title: "Pics",
				URL:   "https://anto.ph",
			},
		},
		Colors: []template.CSS{
			"rgb(245 179 255)",
			"var(--color-white)",
			"var(--color-lime-200)",
			"var(--color-amber-300)",
			"var(--color-blue-200)",
			"var(--color-orange-400)",
		},
		StaticFileHash: staticFileHashes,
	}
}
