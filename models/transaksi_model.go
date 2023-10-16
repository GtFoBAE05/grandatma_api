package models

import "time"

type Transaksi struct {
	Id               int       `json:"id" db:"id"`
	IdReservasi      int       `json:"id_reservasi" db:"id_reservasi"`
	TanggalTransaksi string    `json:"tanggal_transaksi" db:"tanggal_transaksi"`
	TotalPembayaran  float64   `json:"total_pembayaran" db:"total_pembayaran"`
	StatusDeposit    bool      `json:"status_deposit" db:"status_deposit"`
	StatusBayar      bool      `json:"status_bayar" db:"status_bayar"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewTransaksi(idReservasi int, tanggalTransaksi string, totalPembayaran float64, statusDeposit, statusBayar bool) Transaksi {
	return Transaksi{
		IdReservasi:      idReservasi,
		TanggalTransaksi: tanggalTransaksi,
		TotalPembayaran:  totalPembayaran,
		StatusDeposit:    statusDeposit,
		StatusBayar:      statusBayar,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
