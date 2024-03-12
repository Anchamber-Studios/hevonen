package types

type ClubCreate struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Website     string `json:"website" form:"website"`
	Email       string `json:"email" form:"email"`
	Phone       string `json:"phone" form:"phone"`
}

type Club struct {
	ID          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Website     string `json:"website" form:"website"`
	Email       string `json:"email" form:"email"`
	Phone       string `json:"phone" form:"phone"`
}

type ClubMember struct {
	ID    string   `json:"id" form:"id"`
	Name  string   `json:"name" form:"name"`
	Roles []string `json:"roles" form:"roles"`
}

type Member struct {
	ID         string `json:"id" form:"id"`
	ClubID     string `json:"clubId" form:"clubId"`
	FirstName  string `json:"firstName" form:"firstName"`
	MiddleName string `json:"middleName" form:"middleName"`
	LastName   string `json:"lastName" form:"lastName"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
	Height     int    `json:"height" form:"height"`
	Weight     int    `json:"weight" form:"weight"`
}

type MemberCreate struct {
	ClubID     string `json:"clubId" form:"clubId"`
	IdentityID string `json:"identityId" form:"identityId"`
	Email      string `json:"email" form:"email"`
}
