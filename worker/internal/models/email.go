package models

type Email struct {
	Receiver    string `json:"receiver"`
	Topic       string `json:"title"`
	MessageBody string `json:"message_body"`
	ImageUrl    string `json:"image_url"`
}
