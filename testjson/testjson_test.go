package testjson

import "testing"

func TestExampleUsage(t *testing.T) {
	data := []byte(`{
        "title": "Buy milk",
        "assignee": "John Doe",
        "description": "Lorem ipsum dolor sit amet",
		"meta": {
			"meta": {
				"version": 1
			},
			"status": "pending",
			"updated": "1970-01-01T00:00:00Z"
		}
    }`)

	j := Unmarshal(t, data)
	j.Get("title").RequireString(t, "Buy milk")
	j.Get("assignee").RequireString(t, "John Doe")
	j.Get("description").RequireString(t, "Lorem ipsum dolor sit amet")
	j.Get("meta.meta.version").RequireNumber(t, 1)
	j.Get("meta.status").RequireString(t, "pending")
	j.Get("meta.updated").RequireString(t, "1970-01-01T00:00:00Z")
	j.RequireNoAdditionalFields(t)
}
