package domain

type Customer struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Projects []Project `json:"projects"` // has many projects
	UserId   string    `json:"userId"`   // belongs to user
}
