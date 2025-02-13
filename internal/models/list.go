package models

type List struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	ListName      string `json:"list_name"`
	Description   string `json:"description"`
	SharingId     string `json:"sharing_id"`
	ImageFileName string `json:"image_file_name"`
}
