package client

type Member struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Height    int    `json:"height,omitempty"`
	Weight    int    `json:"weight,omitempty"`
}

type MemberCreate struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	Height     int    `json:"height,omitempty"`
}
