package models

import "time"

type StatusMenginap struct {
	Id              int       `json:"id" db:"id"`
	IdReservasi     string    `db:"id_reservasi"`
	Nama            string    `db:"nama"`
	TanggalCheckin  string    `json:"tanggal_checkin" db:"tanggal_checkin"`
	TanggalCheckout string    `json:"tanggal_checkout" db:"tanggal_checkout"`
	StatusCheckin   bool      `json:"status_checkin" db:"status_checkin"`
	StatusCheckout  bool      `json:"status_checkout" db:"status_checkout"`
	TotalPembayaran float64   `db:"total_pembayaran"`
	Jaminan         float64   `json:"jaminan" db:"jaminan"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type StatusComplete struct {
	Id              int       `json:"id" db:"id"`
	IdReservasi     string    `db:"id_reservasi"`
	Nama            string    `db:"nama"`
	TanggalCheckin  string    `json:"tanggal_checkin" db:"tanggal_checkin"`
	TanggalCheckout string    `json:"tanggal_checkout" db:"tanggal_checkout"`
	StatusCheckin   bool      `json:"status_checkin" db:"status_checkin"`
	StatusCheckout  bool      `json:"status_checkout" db:"status_checkout"`
	TotalPembayaran float64   `db:"total_pembayaran"`
	Jaminan         float64   `json:"jaminan" db:"jaminan"`
	Deposit         float64   `json:"deposit" db:"deposit"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
