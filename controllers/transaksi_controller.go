package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransaksiHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Tidak ada token atau Gagal mengambil ID pengguna.",
		})
		return
	}

	var transaksi []models.TransaksiHistory

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi
		FROM
			transaksi t
		JOIN
			reservasi r
		ON 
			t.id_reservasi = r.id_reservasi
		JOIN
			pengguna p
		ON
			r.id_pengguna = p.id
		WHERE
			p.id = $1
	`
	err := database.DBClient.Select(&transaksi, query, userID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  transaksi,
	})
}

func GetTransaksiHistoryByUserId(c *gin.Context) {

	userID := c.Param("userId")

	var transaksi []models.TransaksiHistory

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi
		FROM
			transaksi t
		JOIN
			reservasi r
		ON 
			t.id_reservasi = r.id_reservasi
		JOIN
			pengguna p
		ON
			r.id_pengguna = p.id
		WHERE
			p.id = $1
	`
	err := database.DBClient.Select(&transaksi, query, userID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  transaksi,
	})
}

func GetTransaksiDetail(c *gin.Context) {
	transactionId := c.Param("id")

	var transaksi models.TransaksiDetail
	var fasilitasReservasi []models.FasilitasReservasiXTipeFasilitas

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin, r.tanggal_checkout,
			r.nomor_kamar, r.jumlah_dewasa, r.jumlah_anak, r.nomor_rekening, r.pilihan_kasur
		FROM
			transaksi t
		JOIN
			reservasi r
		ON 
			t.id_reservasi = r.id_reservasi
		WHERE
			t.id_reservasi = $1
	`
	err := database.DBClient.Get(&transaksi, query, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query = `
		SELECT 
			fr.id, fr.id_reservasi, fb.nama_fasilitas, fr.jumlah_unit
		FROM
			transaksi t
		JOIN
			fasilitas_reservasi fr
		ON 
			t.id_reservasi = fr.id_reservasi
		JOIN
			fasilitas_berbayar fb
		ON
			fr.id_fasilitas_berbayar = fb.id
		WHERE
			t.id_reservasi = $1
	`
	err = database.DBClient.Select(&fasilitasReservasi, query, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":     false,
		"data":      transaksi,
		"fasilitas": fasilitasReservasi,
	})
}
