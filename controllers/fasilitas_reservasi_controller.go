package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateFasilitasReservasi(c *gin.Context) {
	var reqBody models.FasilitasReservasi

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	fasilitasReservasi := models.NewFasilitasReservasi(reqBody.IdReservasi, reqBody.IdFasilitasBerbayar, reqBody.JumlahUnit)

	query := `
	INSERT INTO fasilitas_reservasi (id_reservasi, id_fasilitas_berbayar, jumlah_unit, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
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

	_, err = stmt.Exec(fasilitasReservasi.IdReservasi, fasilitasReservasi.IdFasilitasBerbayar, fasilitasReservasi.JumlahUnit, fasilitasReservasi.CreatedAt, fasilitasReservasi.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get total pembayaran sebelumnya
	var tarifSebelum float64
	query = `
		SELECT 
			total_pembayaran
		FROM 
			transaksi 
		WHERE 
			id_reservasi = $1
	`
	err = database.DBClient.Get(&tarifSebelum, query, reqBody.IdReservasi)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get tarif
	var harga float64
	query = `
		SELECT 
			harga
		FROM 
			fasilitas_berbayar 
		WHERE 
			id = $1
	`

	err = database.DBClient.Get(&harga, query, reqBody.IdFasilitasBerbayar)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
			"a":       "a",
		})
		return
	}

	harga = harga * float64(reqBody.JumlahUnit)
	totalPembaranBaru := harga + tarifSebelum

	//update transaksi price
	query = `
		UPDATE transaksi
		SET  total_pembayaran = $1 
		WHERE id_reservasi = $2
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

	_, err = stmt.Exec(totalPembaranBaru, reqBody.IdReservasi)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "SUccess create fasilitas reservasi",
	})
}

func GetFasilitasReservasiByIdReservasi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var fasilitasReservasis []models.FasilitasReservasiXTipeFasilitas
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		SELECT 
			fr.id, fr.id_reservasi, fb.nama_fasilitas, fr.jumlah_unit
			, fr.created_at, fr.updated_at
		FROM
			fasilitas_reservasi fr
		JOIN 
			fasilitas_berbayar fb
		ON 
			fr.id_fasilitas_berbayar = fb.id
		WHERE
			fr.id_reservasi = $1
	`

	err = database.DBClient.Select(&fasilitasReservasis, query, id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  fasilitasReservasis,
	})
}
