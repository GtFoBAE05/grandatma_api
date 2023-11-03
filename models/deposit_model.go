package models

import "time"

type Deposit struct {
	Id          int       `json:"id" db:"id"`
	IdReservasi string    `json:"id_reservasi" db:"id_reservasi"`
	Nominal     float64   `json:"nominal" db:"nominal"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
