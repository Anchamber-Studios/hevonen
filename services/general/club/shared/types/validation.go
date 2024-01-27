package types

import "github.com/anchamber-studios/hevonen/lib"

var (
	ErrClubNameTooShort lib.ValidationError = lib.ValidationError{
		ErrorCode:      "club_name_too_short",
		Message:        "club name must be at least 3 characters long",
		TranslationKey: "club.form.error.name.too_short",
		Field:          "name",
	}
)

func ValidateClubCreate(club ClubCreate) *lib.ValidationError {
	errs := map[string]*lib.ValidationError{}
	if err := ValidateClubName(club.Name); err != nil {
		errs["name"] = err
	}
	return lib.NewValidationErrorForFields(errs)
}

func ValidateClubName(name string) *lib.ValidationError {
	if len(name) < 3 {
		return &ErrClubNameTooShort
	}
	return nil
}
