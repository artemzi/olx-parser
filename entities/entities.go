package entities

import "time"

type DetailsItem struct {
	Name  string `json "name" bson: "name"`
	Value string `json: "value" bson: "value"`
}

type Adverts struct {
	Id        string         `json: "id" bson: "_id"`
	URL       string         `json: "url" bson: "url"`
	Title     string         `json: "title" bson: "title"`
	Place     string         `json: "place" bson: "place"`
	CreatedAt time.Time      `json: "created_at" bson: "created_at"`
	Details   []*DetailsItem `json: "details" bson: "details"`
	Images    []string       `json: "images" bson: "images"`
	Text      string         `json: "text" bson: "text"`
	Price     string         `json: "price" bson: "price"`
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
