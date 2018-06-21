package entities

import "fmt"

type Adverts struct {
	URL    string   `json: "url"`
	Title  string   `json: "title"`
	Place  string   `json: "place"`
	Meta   string   `json: "meta"`
	Images []string `json: "images"`
	Text   string   `json: "text"`
}

func (a Adverts) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
		a.Title, a.Meta, a.Place, a.URL, a.Text, a.Images)
}
