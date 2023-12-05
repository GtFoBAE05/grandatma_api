package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchTransactionAbleToCheckin(c *gin.Context) {
	nama := c.Query("nama")
	id := c.Query("id")

	var transaksi []models.StatusMenginap
	var filteredTransaksi []models.StatusMenginap

	query := `
		SELECT 
			sm.id as id, p.nama as nama, t.id_reservasi, r.tanggal_checkin as tanggal_checkin, r.tanggal_checkout as tanggal_checkout, 
			sm.status_checkin as status_checkin, sm.status_checkout as status_checkout, t.total_pembayaran as total_pembayaran, j.nominal as jaminan
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
		JOIN 
			status_menginap sm
		ON
			t.id_reservasi = sm.id_reservasi
		JOIN 
			jaminan j
		ON t.id_reservasi = j.id_reservasi
		WHERE
			p.nama LIKE '%' || $1 || '%'
		AND
			t.id_reservasi LIKE '%' || $2 || '%'
		AND
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
		AND
			sm.status_checkin = false
		AND
			sm.status_checkout = false
	`
	err := database.DBClient.Select(&transaksi, query, nama, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	for _, element := range transaksi {
		minimumJaminan := element.TotalPembayaran / 2
		if element.Jaminan >= minimumJaminan {
			filteredTransaksi = append(filteredTransaksi, element)
		}
	}

	if len(filteredTransaksi) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Data tidak ada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  filteredTransaksi,
	})

}

func SearchTransactionAbleToCheckout(c *gin.Context) {
	nama := c.Query("nama")
	id := c.Query("id")

	var transaksi []models.StatusMenginap

	query := `
		SELECT 
			sm.id as id, p.nama as nama, t.id_reservasi, r.tanggal_checkin as tanggal_checkin, r.tanggal_checkout as tanggal_checkout, 
			sm.status_checkin as status_checkin, sm.status_checkout as status_checkout
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
		JOIN 
			status_menginap sm
		ON
			t.id_reservasi = sm.id_reservasi
		WHERE
			p.nama LIKE '%' || $1 || '%'
		AND
			t.id_reservasi LIKE '%' || $2 || '%'
		AND
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
		AND
			sm.status_checkin = true
		AND
			sm.status_checkout = false
	`
	err := database.DBClient.Select(&transaksi, query, nama, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if len(transaksi) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Data tidak ada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  transaksi,
	})

}

func SearchTransactionAbleToComplete(c *gin.Context) {
	nama := c.Query("nama")
	id := c.Query("id")

	var transaksi []models.StatusComplete

	query := `
		SELECT 
			sm.id as id, p.nama as nama, t.id_reservasi, r.tanggal_checkin as tanggal_checkin, 
			r.tanggal_checkout as tanggal_checkout, sm.status_checkin as status_checkin, 
			sm.status_checkout as status_checkout, t.total_pembayaran as total_pembayaran,
			d.nominal as deposit, j.nominal as jaminan 

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
		JOIN 
			status_menginap sm
		ON
			t.id_reservasi = sm.id_reservasi
		JOIN
			jaminan j
		on 
			t.id_reservasi = j.id_reservasi
		JOIN
			deposit d
		on 
			t.id_reservasi = d.id_reservasi
		WHERE
			p.nama LIKE '%' || $1 || '%'
		AND
			t.id_reservasi LIKE '%' || $2 || '%'
		AND
			t.status_batal = false
		AND
			sm.status_checkin = true
		AND
			sm.status_checkout = false
		AND 
			t.status_lunas = false
	`
	err := database.DBClient.Select(&transaksi, query, nama, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if len(transaksi) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Data tidak ada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  transaksi,
	})

}

func UpdateStatusCheckin(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Tidak ada id reservasi",
		})
		return
	}

	query := `
		UPDATE status_menginap
		SET  status_checkin = $1, 
			updated_at = $2
		WHERE id_reservasi = $3
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

	_, err = stmt.Exec(true, time.Now(), idStr)

	query = `
		UPDATE deposit
		SET  nominal = $1, 
			updated_at = $2
		WHERE id_reservasi = $3
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

	_, err = stmt.Exec(300000, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated status_checkin",
	})
}

func UpdateStatusCheckout(c *gin.Context) {
	idStr := c.Param("id")

	query := `
		UPDATE status_menginap
		SET  status_checkout = $1, 
			updated_at = $2
		WHERE id_reservasi = $3
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

	_, err = stmt.Exec(true, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query = `
		UPDATE transaksi
		SET  status_lunas = $1, 
			updated_at = $2
		WHERE id_reservasi = $3
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

	_, err = stmt.Exec(true, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated status_checkout",
	})
}
