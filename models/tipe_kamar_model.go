package models

import "time"

type TipeKamar struct {
	Id                 int       `json:"id" db:"id"`
	NamaTipe           string    `json:"nama_tipe" db:"nama_tipe"`
	PilihanTempatTidur string    `json:"pilihan_tempat_tidur" db:"pilihan_tempat_tidur"`
	Fasilitas          string    `json:"fasilitas" db:"fasilitas"`
	Deskripsi          string    `json:"deskripsi" db:"deskripsi"`
	RincianKamar       string    `json:"rincian_kamar" db:"rincian_kamar"`
	Tarif              float64   `json:"tarif" db:"tarif"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

func NewTipeKamar(namaTipe, pilihanTempatTidur, fasilitas, deskripsi, rincianKamar string, tarif float64) TipeKamar {
	return TipeKamar{
		NamaTipe:           namaTipe,
		PilihanTempatTidur: pilihanTempatTidur,
		Fasilitas:          fasilitas,
		Deskripsi:          deskripsi,
		RincianKamar:       rincianKamar,
		Tarif:              tarif,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
