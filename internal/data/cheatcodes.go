package data

import (
	"time"

	"github.com/arnokay/gta-vc-cheatcodes-crud/internal/validator"
)

type Cheatcode struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Code        *string   `json:"code"`
	Description *string   `json:"description"`
	Tags        []string  `json:"tags"`
	Version     int32     `json:"version"`
}

func ValidateCheatcode(v *validator.Validator, cheatcode *Cheatcode) {
	v.Check(*cheatcode.Code != "", "code", "must be provided")
	v.Check(len(*cheatcode.Code) <= 500, "code", "must not be more than 500 bytes long")

	if cheatcode.Description != nil {
		v.Check(*cheatcode.Description != "", "description", "must not be empty")
		v.Check(len(*cheatcode.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	}

	v.Check(cheatcode.Tags != nil, "tags", "must be provided")
	v.Check(len(cheatcode.Tags) >= 1, "tags", "must contain at least 1 tag")
	v.Check(len(cheatcode.Tags) <= 5, "tags", "must not contain more than 5 tags")
	v.Check(validator.Unique(cheatcode.Tags), "tags", "must not contain duplicate values")
}
