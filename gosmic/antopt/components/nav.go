package components

import "html/template"

type NavItem struct {
	Title string
	URL   string
}

type Nav struct {
	Nav         []NavItem
	CurrentURL  string
	ColorPicker ColorPicker
}

func NewNav(items []NavItem, currentURL string, colors []template.CSS, selectedColor template.CSS) Nav {
	return Nav{
		Nav:        items,
		CurrentURL: currentURL,
		ColorPicker: ColorPicker{
			Colors:   colors,
			Selected: selectedColor,
		},
	}
}
