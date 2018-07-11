package domain

type User struct {
	Id       string
	Email    string `json:"email"`
	Password string `json:"password"`
}
