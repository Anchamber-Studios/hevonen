package client

type UserCreate struct {
	Email    string          `json:"email" form:"email"`
	Username string          `json:"username" form:"username"`
	Password string          `json:"password" form:"password"`
	Apps     []AppConnection `json:"apps" form:"apps"`
}

type AppConnection struct {
	App   string `json:"app"`
	Token string `json:"token"`
}

type User struct {
	Id       string   `json:"id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Apps     []string `json:"apps"`
}
