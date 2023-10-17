package models

import "time"

type Transaksi struct {
	Id               int       `json:"id" db:"id"`
	IdReservasi      string    `json:"id_reservasi" db:"id_reservasi"`
	TanggalTransaksi string    `json:"tanggal_transaksi" db:"tanggal_transaksi"`
	TotalPembayaran  float64   `json:"total_pembayaran" db:"total_pembayaran"`
	StatusDeposit    bool      `json:"status_deposit" db:"status_deposit"`
	StatusBayar      bool      `json:"status_bayar" db:"status_bayar"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type TransaksiHistory struct {
	IdReservasi      string  `json:"id_reservasi" db:"id_reservasi"`
	TanggalTransaksi string  `json:"tanggal_transaksi" db:"tanggal_transaksi"`
	TotalPembayaran  float64 `json:"total_pembayaran" db:"total_pembayaran"`
}

type TransaksiDetail struct {
	IdReservasi      string  `db:"id_reservasi"`
	TanggalTransaksi string  `db:"tanggal_transaksi"`
	TotalPembayaran  float64 `db:"total_pembayaran"`
	NomorKamar       int     `db:"nomor_kamar"`
	TanggalCheckin   string  `db:"tanggal_checkin"`
	TanggalCheckout  string  `db:"tanggal_checkout"`
	JumlahDewasa     int     `db:"jumlah_dewasa"`
	JumlahAnak       int     `db:"jumlah_anak"`
	NomorRekening    string  `db:"nomor_rekening"`
	PilihanKasur     string  `db:"pilihan_kasur"`
}

func NewTransaksi(idReservasi string, tanggalTransaksi string, totalPembayaran float64, statusDeposit, statusBayar bool) Transaksi {
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
