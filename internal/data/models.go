package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Cheatcodes interface {
		Insert(cheatcode *Cheatcode) error
		Get(id int64) (*Cheatcode, error)
		Update(cheatcode *Cheatcode) error
		Delete(id int64) error
		GetAll(code string, description string, tags []string, filters Filters) ([]*Cheatcode, Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cheatcodes: CheatcodeModel{DB: db},
		Users:      UserModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Cheatcodes: MockCheatcodeModel{},
	}
}
