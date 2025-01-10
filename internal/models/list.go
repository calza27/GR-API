package models

type List struct {
	Uuid      string `json:"uuid"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	SharingId string `json:"sharing_id"`
}
