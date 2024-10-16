package entity

import "github.com/google/uuid"

type Doctor struct {
	Id      string
	Name    string
	Surname string
	Cabinet string
	Type    string
}

func NewDoctor(name, surname, cabinet, typ string) *Doctor {
	return &Doctor{
		Id:      uuid.NewString(),
		Name:    name,
		Surname: surname,
		Cabinet: cabinet,
		Type:    typ,
	}
}
