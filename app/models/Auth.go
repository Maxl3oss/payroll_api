package models

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ReqToken struct {
	Refresh string `json:"refresh"`
}
