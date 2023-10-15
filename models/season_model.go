package models

import "time"

type Season struct {
	Id              int       `json:"id" db:"id"`
	NamaSeason      string    `json:"nama_season" db:"nama_season"`
	TanggalMulai    string    `json:"tanggal_mulai" db:"tanggal_mulai"`
	TanggalBerakhir string    `json:"tanggal_berakhir" db:"tanggal_berakhir"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewSeason(namaSeason, tanggalMulai, tanggalBerakhir string) Season {
	return Season{
		NamaSeason:      namaSeason,
		TanggalMulai:    tanggalMulai,
		TanggalBerakhir: tanggalBerakhir,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
