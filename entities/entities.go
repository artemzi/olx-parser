package entities

import "time"

type DetailsItem struct {
	Name  string `json "name"`
	Value string `json: "value"`
}

type Adverts struct {
	Id        string         `json: "id"`
	URL       string         `json: "url"`
	Title     string         `json: "title"`
	Place     string         `json: "place"`
	CreatedAt time.Time      `json: "created_at"`
	Details   []*DetailsItem `json: "details"`
	Images    []string       `json: "images"`
	Text      string         `json: "text"`
	Price     string         `json: "price"`
	Phone     string         `json: "phone,omitempty"`
}

// TODO: do i need it?
//func (a Adverts) String() string {
//	return fmt.Sprintf("%s\n%v\n%s\n%s\n%s\n%s\n%s\n%s\n",
//		a.Title, a.Details, a.Meta, a.Place, a.URL, a.Text, a.Images, a.Price)
//}

type AdvertsResponse struct {
	Size    int        `json: "count"`
	Adverts []*Adverts `json: "data"`
}
