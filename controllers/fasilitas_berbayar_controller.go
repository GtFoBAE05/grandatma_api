package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateFasilitasBerbayar(c *gin.Context) {
	var reqBody models.FasilitasBerbayar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	fasilitasBerbayar := models.NewFasilitasBerbayar(reqBody.NamaFasilitas, reqBody.Harga)

	query := `
	INSERT INTO fasilitas_berbayar (nama_fasilitas, harga, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
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

	_, err = stmt.Exec(fasilitasBerbayar.NamaFasilitas, fasilitasBerbayar.Harga, fasilitasBerbayar.CreatedAt, fasilitasBerbayar.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success create fasilitas berbayar",
	})
}

func GetFasilitasBerbayars(c *gin.Context) {
	var fasilitasBerbayars []models.FasilitasBerbayar

	query := `
		SELECT 
			id, nama_fasilitas, harga, created_at, updated_at
		FROM
			fasilitas_berbayar
	`

	err := database.DBClient.Select(&fasilitasBerbayars, query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  fasilitasBerbayars,
	})
}

func GetFasilitasBerbayarById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var fasilitasBerbayar models.FasilitasBerbayar
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		SELECT 
			id, nama_fasilitas, harga, created_at, updated_at
		FROM
			fasilitas_berbayar
		WHERE id = $1
	`

	err = database.DBClient.Get(&fasilitasBerbayar, query, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  fasilitasBerbayar,
	})

}

func UpdateFasilitasBerbayar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var reqBody models.FasilitasBerbayar

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE fasilitas_berbayar
		SET  nama_fasilitas = $1, 
			harga = $2, 
			 updated_at = $3
		WHERE id = $4
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

	_, err = stmt.Exec(reqBody.NamaFasilitas, reqBody.Harga, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated fasilitas berbayar",
	})

}

func DeleteFasilitasBerbayar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	res, err := database.DBClient.Exec("DELETE FROM fasilitas_berbayar WHERE id = $1", id)

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
			"message": "Gagal Menghapus, id fasilitas berbayar tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Berhasil hapus fasilitas berbayar",
	})
}
