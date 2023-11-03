package models

import "time"

type Reservasi struct {
	Id              int       `json:"id" db:"id"`
	IdReservasi     string    `json:"id_reservasi" db:"id_reservasi"`
	Email_pengguna  string    `json:"email_pengguna" db:"email_pengguna"`
	IdKamar         int       `json:"id_kamar" db:"id_kamar"`
	TanggalCheckin  string    `json:"tanggal_checkin" db:"tanggal_checkin"`
	TanggalCheckout string    `json:"tanggal_checkout" db:"tanggal_checkout"`
	JumlahDewasa    int       `json:"jumlah_dewasa" db:"jumlah_dewasa"`
	JumlahAnak      int       `json:"jumlah_anak" db:"jumlah_anak"`
	NomorRekening   string    `json:"nomor_rekening" db:"nomor_rekening"`
	PilihanKasur    string    `json:"pilihan_kasur" db:"pilihan_kasur"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type TotalReservasi struct {
	TotalReservasi int `db:"total_reservasi"`
}

func NewReservasi(idReservasi string, emailPengguna string, idKamar int, tanggalCheckin string, tanggalCheckout string, jumlahDewasa int, jumlahAnak int, nomorRekening string, pilihanKasur string) Reservasi {
	return Reservasi{
		IdReservasi:     idReservasi,
		Email_pengguna:  emailPengguna,
		IdKamar:         idKamar,
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
