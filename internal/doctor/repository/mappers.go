package repository

import (
	"github.com/neiasit/service-boilerplate/internal/doctor/entity"
	userEntity "github.com/neiasit/service-boilerplate/internal/user/entity"
)

func MapDBToDomainDoctor(doctorDB *DoctorDB) *entity.Doctor {
	return &entity.Doctor{
		Id:      doctorDB.Id,
		Name:    doctorDB.Name,
		Surname: doctorDB.Surname,
		Cabinet: doctorDB.Cabinet,
		Type:    doctorDB.Type,
	}
}

func MapDomainToDBDoctor(doctor *entity.Doctor) *DoctorDB {
	return &DoctorDB{
		Id:      doctor.Id,
		Name:    doctor.Name,
		Surname: doctor.Surname,
		Cabinet: doctor.Cabinet,
		Type:    doctor.Type,
	}
}

func MapDBToDomainUser(userDB *UserDB) *userEntity.User {
	return &userEntity.User{
		Id:       userDB.Id,
		Name:     userDB.Name,
		Surname:  userDB.Surname,
		Verdict:  userDB.Verdict,
		IsLoсked: userDB.IsLocked,
	}
}

func MapDomainToDBUser(user *userEntity.User) *UserDB {
	return &UserDB{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		IsLocked: user.IsLoсked,
		Verdict:  user.Verdict,
	}
}
