package models

import "time"

type Reservasi struct {
	Id              int       `json:"id" db:"id"`
	IdReservasi     string    `json:"id_reservasi" db:"id_reservasi"`
	IdPengguna      int       `json:"id_pengguna" db:"id_pengguna"`
	NomorKamar      int       `json:"nomor_kamar" db:"nomor_kamar"`
	TanggalCheckin  string    `json:"tanggal_checkin" db:"tanggal_checkin"`
	TanggalCheckout string    `json:"tanggal_checkout" db:"tanggal_checkout"`
	JumlahDewasa    int       `json:"jumlah_dewasa" db:"jumlah_dewasa"`
	JumlahAnak      int       `json:"jumlah_anak" db:"jumlah_anak"`
	NomorRekening   string    `json:"nomor_rekening" db:"nomor_rekening"`
	PilihanKasur    string    `json:"pilihan_kasur" db:"pilihan_kasur"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewReservasi(idReservasi string, idPengguna int, nomorKamar int, tanggalCheckin string, tanggalCheckout string, jumlahDewasa int, jumlahAnak int, nomorRekening string, pilihanKasur string) Reservasi {
	return Reservasi{
		IdReservasi:     idReservasi,
		IdPengguna:      idPengguna,
		NomorKamar:      nomorKamar,
		TanggalCheckin:  tanggalCheckin,
		TanggalCheckout: tanggalCheckout,
		JumlahDewasa:    jumlahDewasa,
		JumlahAnak:      jumlahAnak,
		NomorRekening:   nomorRekening,
		PilihanKasur:    pilihanKasur,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
