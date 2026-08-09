package entity

type User struct {
	ID          uint `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password	string `json:"-"`
}
