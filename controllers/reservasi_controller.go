package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateReservasi(c *gin.Context) {
	var reqBody models.Reservasi

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	reservasi := models.NewReservasi(reqBody.IdReservasi, reqBody.IdPengguna, reqBody.NomorKamar, reqBody.TanggalCheckin, reqBody.TanggalCheckout, reqBody.JumlahDewasa, reqBody.JumlahAnak, reqBody.NomorRekening, reqBody.PilihanKasur)

	//add reservasi
	query := `
	INSERT INTO reservasi (id_reservasi, id_pengguna, nomor_kamar, tanggal_checkin, tanggal_checkout, 
		jumlah_dewasa, jumlah_anak, nomor_rekening, pilihan_kasur, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	stmt, err := database.DBClient.Prepare(query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(reservasi.IdReservasi, reservasi.IdPengguna, reservasi.NomorKamar, reservasi.TanggalCheckin, reservasi.TanggalCheckout,
		reservasi.JumlahDewasa, reservasi.JumlahAnak, reservasi.NomorRekening, reservasi.PilihanKasur, reservasi.CreatedAt, reservasi.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get tarif kamar
	var tarifKamar models.TarifKamar
	query = `
		SELECT 
			t.tarif as tarif
		FROM 
			kamar k
		JOIN 
			tarif t
		ON
			k.id_tipe_kamar = t.id_tipe_kamar
		WHERE 
			k.nomor_kamar = $1
	`

	err = database.DBClient.Get(&tarifKamar, query, reservasi.NomorKamar)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
			"error2":  "no rows",
		})
		return
	}

	//get durasi kamar
	layout := "2006-01-02"

	tanggalMulai, err := time.Parse(layout, reservasi.TanggalCheckin)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	tanggalSelesai, err := time.Parse(layout, reservasi.TanggalCheckout)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	diff := tanggalSelesai.Sub(tanggalMulai)
	diffInDays := int(diff.Hours() / 24)

	//total pembayaran
	totalPembayaran := float64(diffInDays) * tarifKamar.Tarif

	//add transaksi
	query = `
	INSERT INTO transaksi (id_reservasi,tanggal_transaksi, total_pembayaran, status_deposit, status_bayar, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	stmt, err = database.DBClient.Prepare(query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(reservasi.IdReservasi, time.Now(), totalPembayaran, false, false, reservasi.CreatedAt, reservasi.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success create reservasi",
	})
}
