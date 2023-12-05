package models

import "time"

type Pengguna struct {
	Id        int       `json:"id" db:"id"`
	Nama      string    `json:"nama" db:"nama"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Notelp    string    `json:"notelp" db:"notelp"`
	Password  string    `json:"password" db:"password"`
	Alamat    string    `json:"alamat" db:"alamat"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Register struct {
	Nama     string `json:"nama" db:"nama"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Alamat   string `json:"alamat" db:"alamat"`
	Notelp   string `json:"notelp" db:"notelp"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type ChangePassword struct {
	Password string `json:"password" db:"password"`
}

type UpdateProfile struct {
	Nama     string `json:"nama" db:"nama"`
	Email    string `json:"email" db:"email"`
	Alamat   string `json:"alamat" db:"alamat"`
	Username string `json:"username" db:"username"`
	Notelp   string `json:"notelp" db:"notelp"`
}

func NewPengguna(nama, email, username, notelp, password, alamat, role string) Pengguna {
	return Pengguna{
		Nama:      nama,
		Email:     email,
		Username:  username,
		Notelp:    notelp,
		Password:  password,
		Alamat:    alamat,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (p Pengguna) WithId(id int) Pengguna {
	p.Id = id
	return p
}
