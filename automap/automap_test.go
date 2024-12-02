package automap_test

import (
	"fmt"

	"anto.pt/x/automap"
)

type Post struct {
	Title  string
	Body   string
	Secret string
}

func (p Post) String() string {
	return fmt.Sprintf("Post{title=%s, body=%s}", p.Title, p.Body)
}

type PostDTO struct {
	Title string
	Body  string
}

func (d PostDTO) String() string {
	return fmt.Sprintf("PostDTO{title=%s, body=%s}", d.Title, d.Body)
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
	// Output: PostDTO{title=The title, body=Did you see this?}
}
