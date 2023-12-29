package client

import "time"

type UserCreate struct {
	Email    string          `json:"email" form:"email"`
	Password string          `json:"password" form:"password"`
	Apps     []AppConnection `json:"apps" form:"apps"`
}

type AppConnection struct {
	App   string `json:"app"`
	Token string `json:"token"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	EmailConfirmed bool      `json:"emailConfirmed"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	Apps           []string  `json:"apps"`
}
