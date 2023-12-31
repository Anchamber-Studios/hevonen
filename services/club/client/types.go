package client

type ClubCreate struct {
	Name    string `json:"name" form:"name"`
	Website string `json:"website" form:"website"`
}

type Member struct {
	ID         string `json:"id" form:"id"`
	ClubID     uint64 `json:"clubId" form:"clubId"`
	FirstName  string `json:"firstName" form:"firstName"`
	MiddleName string `json:"middleName" form:"middleName"`
	LastName   string `json:"lastName" form:"lastName"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
	Height     int    `json:"height" form:"height"`
	Weight     int    `json:"weight" form:"weight"`
}

type MemberCreate struct {
	ClubID     uint64 `json:"clubId" form:"clubId"`
	FirstName  string `json:"firstName" form:"firstName"`
	MiddleName string `json:"middleName" form:"middleName"`
	LastName   string `json:"lastName" form:"lastName"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
	Weight     int    `json:"weight" form:"weight"`
	Height     int    `json:"height" form:"height"`
}
