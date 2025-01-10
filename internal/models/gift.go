package models

type Gift struct {
	Uuid            string `json:"uuid"`
	ListId          string `json:"list_id"`
	CreatedAt       string `json:"created_at"`
	Title           string `json:"title"`
	Desription      string `json:"description"`
	PlaceOfPurchase string `json:"place_of_purchase"`
	ImageFileName   string `json:"image_file_name"`
	Url             string `json:"url"`
	Price           int    `json:"price"`
	Rank            int    `json:"rank"`
}
