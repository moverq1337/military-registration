package repository

type DoctorDB struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	Surname string `db:"surname"`
	Cabinet string `db:"cabinet"`
	Type    string `db:"type"`
}

type UserDB struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Verdict  bool   `db:"verdict"`
	IsLocked bool   `db:"is_locked"`
}
