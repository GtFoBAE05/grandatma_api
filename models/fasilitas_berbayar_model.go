package models

import "time"

type FasilitasBerbayar struct {
	Id            int       `json:"id" db:"id"`
	NamaFasilitas string    `json:"nama_fasilitas" db:"nama_fasilitas"`
	Harga         float64   `json:"harga" db:"harga"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewFasilitasBerbayar(namaFasilitas string, harga float64) FasilitasBerbayar {
	return FasilitasBerbayar{
		NamaFasilitas: namaFasilitas,
		Harga:         harga,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
