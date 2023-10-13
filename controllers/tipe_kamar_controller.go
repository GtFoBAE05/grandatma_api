package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTipeKamar(c *gin.Context) {
	var reqBody models.TipeKamar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	tipeKamar := models.NewTipeKamar(reqBody.NamaTipe, reqBody.PilihanTempatTidur, reqBody.Fasilitas, reqBody.Deskripsi, reqBody.RincianKamar)

	query := `
	INSERT INTO tipe_kamar (nama_tipe, pilihan_tempat_tidur, fasilitas, deskripsi, rincian_kamar, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
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

	_, err = stmt.Exec(tipeKamar.NamaTipe, tipeKamar.PilihanTempatTidur, tipeKamar.Fasilitas, tipeKamar.Deskripsi, tipeKamar.RincianKamar, tipeKamar.CreatedAt, tipeKamar.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "SUccess create tipe kamar",
	})
}

func GetTipeKamar(c *gin.Context) {
	var tipeKamars []models.TipeKamar

	query := `
		SELECT 
			id, nama_tipe, pilihan_tempat_tidur
			, fasilitas, deskripsi, rincian_kamar
			, created_at, updated_at
		FROM
			tipe_kamar
	`

	err := database.DBClient.Select(&tipeKamars, query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  tipeKamars,
	})
}

func GetTipeKamarById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var tipeKamar models.TipeKamar
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		SELECT 
			id, nama_tipe, pilihan_tempat_tidur
			, fasilitas, deskripsi, rincian_kamar
			, created_at, updated_at
		FROM
			tipe_kamar
		WHERE id = $1
	`

	err = database.DBClient.Get(&tipeKamar, query, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  tipeKamar,
	})

}

func UpdateTipeKamar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var reqBody models.TipeKamar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE tipe_kamar
		set  nama_tipe = $1, 
			pilihan_tempat_tidur = $2, 
			fasilitas = $3, 
			deskripsi = $4, 
			rincian_kamar = $5,
			 updated_at = $6
		WHERE id = $7
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

	_, err = stmt.Exec(reqBody.NamaTipe, reqBody.PilihanTempatTidur, reqBody.Fasilitas,
		reqBody.Deskripsi, reqBody.RincianKamar, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated tipe kamar",
	})

}

func DeleteTipeKamar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	res, err := database.DBClient.Exec("DELETE FROM tipe_kamar WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Gagal Menghapus, id tipe kamar tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Berhasil hapus menu",
	})
}
