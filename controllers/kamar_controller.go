package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateKamar(c *gin.Context) {
	var reqBody models.Kamar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	kamar := models.NewKamar(reqBody.NomorKamar, reqBody.IdTipeKamar, reqBody.Status)

	query := `
	INSERT INTO kamar (nomor_kamar, id_tipe_kamar, status, created_at, updated_at)
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

	_, err = stmt.Exec(kamar.NomorKamar, kamar.IdTipeKamar, kamar.Status, kamar.CreatedAt, kamar.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "SUccess create kamar",
	})
}

func GetKamar(c *gin.Context) {
	var kamars []models.KamarXTipeKamar

	query := `
		SELECT 
			k.id as id_kamar, k.nomor_kamar, tk.nama_tipe, k.status
			, k.created_at, k.updated_at
		FROM
			kamar k
		JOIN 
			tipe_kamar tk
		on k.id_tipe_kamar = tk.id
	`

	err := database.DBClient.Select(&kamars, query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  kamars,
	})
}

func GetKamarById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var kamar models.Kamar
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		SELECT 
			id, nomor_kamar,id_tipe_kamar, status
			, created_at, updated_at
		FROM
			kamar
		WHERE id = $1
	`

	err = database.DBClient.Get(&kamar, query, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  kamar,
	})

}

func UpdateKamar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var reqBody models.Kamar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE kamar
		SET  nomor_kamar = $1, 
			id_tipe_kamar = $2, 
			status = $3, 
			 updated_at = $4
		WHERE id = $5
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

	_, err = stmt.Exec(reqBody.NomorKamar, reqBody.IdTipeKamar, reqBody.Status, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated kamar",
	})

}

func DeleteKamar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	res, err := database.DBClient.Exec("DELETE FROM kamar WHERE id = $1", id)

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
			"message": "Gagal Menghapus, id kamar tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Berhasil hapus kamar",
	})
}
