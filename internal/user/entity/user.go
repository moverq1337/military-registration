package entity

import "github.com/google/uuid"

type User struct {
	Id       string
	Name     string
	Surname  string
	IsLoсked bool
	Verdict  bool
}

func NewUser(name, surname string) *User {
	return &User{
		Id:       uuid.NewString(),
		Name:     name,
		Surname:  surname,
		IsLoсked: false,
		Verdict:  false,
	}
}
