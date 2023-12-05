package models

type Jaminan struct {
	Id          int64   `json:"id" db:"id"`
	IdReservasi string  `json:"id_reservasi" db:"id_reservasi"`
	Nominal     float64 `json:"nominal" db:"nominal"`
	StatusLunas bool    `json:"status_lunas" db:"status_lunas"`
}

type ListJaminan struct {
	Id              int64   `json:"id" db:"id"`
	IdReservasi     string  `json:"id_reservasi" db:"id_reservasi"`
	Nominal         float64 `json:"nominal" db:"nominal"`
	TotalPembayaran float64 `json:"total_pembayaran" db:"total_pembayaran"`
	StatusLunas     bool    `json:"status_lunas" db:"status_lunas"`
}

type SearchUncompletedJaminan struct {
	Nama             string  `json:"nama" db:"nama"`
	IdReservasi      string  `json:"id_reservasi" db:"id_reservasi"`
	TanggalTransaksi string  `json:"tanggal_transaksi" db:"tanggal_transaksi"`
	TanggalCheckin   string  `json:"tanggal_checkin" db:"tanggal_checkin"`
	TotalJaminan     float64 `json:"total_jaminan" db:"total_jaminan"`
	TotalPembayaran  float64 `json:"total_pembayaran" db:"total_pembayaran"`
	StatusLunas      bool    `json:"status_lunas" db:"status_lunas"`
}

type UpdateJaminanPayload struct {
	Nominal float64 `json:"nominal" db:"nominal"`
}

type UpdateJaminanWithRekeningPayload struct {
	Nominal  float64 `json:"nominal" db:"nominal"`
	Rekening string  `json:"rekening" db:"rekening"`
}
