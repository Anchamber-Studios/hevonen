package client

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
	Id    string   `json:"id"`
	Email string   `json:"email"`
	Apps  []string `json:"apps"`
}
