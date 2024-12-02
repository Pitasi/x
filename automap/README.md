```
package main

import (
	"fmt"

	"anto.pt/x/automap"
)

type Post struct {
	Title  string
	Body   string
	Secret string
}

type PostDTO struct {
	Title string
	Body  string
}

func Example() {
	automap.Register(func(post Post) PostDTO {
		return PostDTO{
			Title: post.Title,
			Body:  post.Body,
		}
	})

	p := Post{
		Title:  "The title",
		Body:   "Did you see this?",
		Secret: "This is a secret",
	}

	var dto PostDTO
	automap.Map(p, &dto)

	fmt.Println(dto)
}
```
