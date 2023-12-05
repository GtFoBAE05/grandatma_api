package models

type NewCustomerStatisticsByYear struct {
	January   int `json:"january" db:"january"`
	February  int `json:"february" db:"february"`
	March     int `json:"march" db:"march"`
	April     int `json:"april" db:"april"`
	May       int `json:"may" db:"may"`
	June      int `json:"june" db:"june"`
	July      int `json:"july" db:"july"`
	August    int `json:"august" db:"august"`
	September int `json:"september" db:"september"`
	October   int `json:"october" db:"october"`
	November  int `json:"november" db:"november"`
	December  int `json:"december" db:"december"`
	Total     int `json:"total" db:"total"`
}

type MonthlyIncomeReport struct {
	Month int    `json:"month" db:"month"`
	Type  string `json:"type" db:"type"`
	Total int    `json:"total" db:"total"`
}

type RoomType struct {
	IDTipeKamar   int    `json:"id_tipe_kamar" db:"id_tipe_kamar"`
	NamaTipeKamar string `json:"nama_tipe_kamar" db:"nama_tipe_kamar"`
}

type RoomVisitorReport struct {
	Type     string `json:"type" db:"type"`
	Personal int    `json:"personal" db:"personal"`
	Group    int    `json:"group" db:"group"`
	Total    int    `json:"total" db:"total"`
}

type TopCustomer struct {
	Name             string `json:"name" db:"name"`
	TotalReservation string `json:"reservation_count" db:"reservation_count"`
	TotalPayment     string `json:"total_payment" db:"total_payment"`
}
