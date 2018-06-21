package entities

import "fmt"

type DetailsItem struct {
	Name  string `json "name"`
	Value string `json: "value"`
}

type Adverts struct {
	URL     string         `json: "url"`
	Title   string         `json: "title"`
	Place   string         `json: "place"`
	Meta    string         `json: "meta"`
	Details []*DetailsItem `json: "details"`
	Images  []string       `json: "images"`
	Text    string         `json: "text"`
	Price   string         `json: "price"`
	Phone   string         `json: "phone,omitempty"`
}

func (a Adverts) String() string {
	return fmt.Sprintf("%s\n%v\n%s\n%s\n%s\n%s\n%s\n%s\n",
		a.Title, a.Details, a.Meta, a.Place, a.URL, a.Text, a.Images, a.Price)
}

type AdvertsResponse struct {
	Size    int        `json: "count"`
	Adverts []*Adverts `json: "data"`
}
