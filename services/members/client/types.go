package client

type Member struct {
	ID         string `json:"id" form:"id"`
	ClubID     int    `json:"clubId" form:"clubId"`
	FirstName  string `json:"firstName" form:"firstName"`
	MiddleName string `json:"middleName" form:"middleName"`
	LastName   string `json:"lastName" form:"lastName"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
	Height     int    `json:"height" form:"height"`
	Weight     int    `json:"weight" form:"weight"`
}

type MemberCreate struct {
	ClubID     int    `json:"clubId" form:"clubId"`
	FirstName  string `json:"firstName" form:"firstName"`
	MiddleName string `json:"middleName" form:"middleName"`
	LastName   string `json:"lastName" form:"lastName"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
	Weight     int    `json:"weight" form:"weight"`
	Height     int    `json:"height" form:"height"`
}
