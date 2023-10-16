package models

import "time"

type FasilitasReservasi struct {
	Id                  int       `json:"id" db:"id"`
	IdReservasi         int       `json:"id_reservasi" db:"id_reservasi"`
	IdFasilitasBerbayar int       `json:"id_fasilitas_berbayar" db:"id_fasilitas_berbayar"`
	JumlahUnit          int       `json:"jumlah_unit" db:"jumlah_unit"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type FasilitasReservasiXTipeFasilitas struct {
	Id            int       `db:"id"`
	IdReservasi   int       `db:"id_reservasi"`
	NamaFasilitas int       `db:"nama_fasilitas"`
	JumlahUnit    int       `db:"jumlah_unit"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewFasilitasReservasi(idReservasi, idFasilitasBerbayar, jumlahUnit int) FasilitasReservasi {
	return FasilitasReservasi{
		IdReservasi:         idReservasi,
		IdFasilitasBerbayar: idFasilitasBerbayar,
		JumlahUnit:          jumlahUnit,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
