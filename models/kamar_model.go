package models

import "time"

type Kamar struct {
	Id          int       `json:"id" db:"id"`
	NomorKamar  string    `json:"nomor_kamar" db:"nomor_kamar"`
	IdTipeKamar int       `json:"id_tipe_kamar" db:"id_tipe_kamar"`
	Status      bool      `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type KamarXTipeKamar struct {
	Id         int       `db:"id_kamar"`
	NomorKamar string    `db:"nomor_kamar"`
	NamaTipe   string    `db:"nama_tipe"`
	Status     bool      `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func NewKamar(nomorKamar string, idTipeKamar int, status bool) Kamar {
	return Kamar{
		NomorKamar:  nomorKamar,
		IdTipeKamar: idTipeKamar,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
