package pages

import (
	"html/template"

	"anto.pt/x/gosmic/antopt/components"
)

type Common struct {
	BaseHeaders components.BaseHeaders
	Nav         components.Nav
	Debug       components.Debug
}

func NewCommon(navItems []components.NavItem, currentURL string, colors []template.CSS, selectedColor template.CSS, dbg components.Debug) Common {
	return Common{
		BaseHeaders: components.BaseHeaders{
			SelectedColor: selectedColor,
		},
		Nav:   components.NewNav(navItems, currentURL, colors, selectedColor),
		Debug: dbg,
	}
}
