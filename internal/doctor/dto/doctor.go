package dto

type CreateDoctorRs struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Cabinet string `json:"cabinet"`
	Type    string `json:"type"`
}

type NextRq struct {
	DoctorId string `json:"doctorId"`
	UserId   string `json:"userId"`
	Verdict  bool   `json:"verdict"`
}

type TherapistNextRq struct {
	UserId  string `json:"userId"`
	Verdict bool   `json:"verdict"`
}
