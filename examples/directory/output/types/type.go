package types

import (
	"time"
)

type (
	ID int

	User struct {
		ID        *ID        `json:"id"`
		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}
)

func (id ID) IsValid() bool {

	if id > 0 {
		return true
	}

	return false
}
