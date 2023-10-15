package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTarif(c *gin.Context) {
	var reqBody models.Tarif

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	tarif := models.NewTarif(reqBody.IdTipeKamar, reqBody.SeasonId, reqBody.Tarif)

	query := `
	INSERT INTO tarif (id_tipe_kamar, season_id, tarif, created_at, updated_at)
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

	_, err = stmt.Exec(tarif.IdTipeKamar, tarif.SeasonId, tarif.Tarif, tarif.CreatedAt, tarif.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success create tarif",
	})
}

func GetTarifs(c *gin.Context) {
	var tarifs []models.TarifXTipeKamarXSeason

	query := `
		SELECT 
			t.id, tk.nama_tipe as nama_tipe_kamar, s.nama_season, t.tarif, t.created_at, t.updated_at
		FROM
			tarif t
		JOIN
			tipe_kamar tk
		ON 
			t.id_tipe_kamar = tk.id
		JOIN
			season s
		ON
			t.season_id = s.id
	`

	err := database.DBClient.Select(&tarifs, query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  tarifs,
	})
}

func GetTarifById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var tarif models.TarifXTipeKamarXSeason
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		SELECT 
			t.id, tk.nama_tipe as nama_tipe_kamar, s.nama_season, t.tarif, t.created_at, t.updated_at
		FROM
			tarif t
		JOIN
			tipe_kamar tk
		ON 
			t.id_tipe_kamar = tk.id
		JOIN
			season s
		ON
			t.season_id = s.id
		WHERE t.id = $1
	`

	err = database.DBClient.Get(&tarif, query, id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  tarif,
	})

}

func UpdateTarif(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var reqBody models.Tarif

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE tarif
		SET  id_tipe_kamar = $1, 
			season_id = $2,
			tarif = $3, 
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

	_, err = stmt.Exec(reqBody.IdTipeKamar, reqBody.SeasonId, reqBody.Tarif, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated tarif",
	})

}

func DeleteTarif(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	res, err := database.DBClient.Exec("DELETE FROM tarif WHERE id = $1", id)

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
			"message": "Gagal Menghapus, id tarif tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "Berhasil hapus tarif",
	})
}
