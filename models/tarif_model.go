package models

import "time"

type Tarif struct {
	Id          int       `json:"id" db:"id"`
	IdTipeKamar int       `json:"id_tipe_kamar" db:"id_tipe_kamar"`
	SeasonId    int       `json:"id_season" db:"season_id"`
	Tarif       float64   `json:"tarif" db:"tarif"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TarifXTipeKamarXSeason struct {
	Id         int       `db:"id"`
	TipeKamar  string    `db:"nama_tipe_kamar"`
	NamaSeason string    `db:"nama_season"`
	Tarif      float64   `db:"tarif"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func NewTarif(idTipeKamar int, seasonId int, tarif float64) Tarif {
	return Tarif{
		IdTipeKamar: idTipeKamar,
		SeasonId:    seasonId,
		Tarif:       tarif,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
