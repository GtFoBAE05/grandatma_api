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

func GetKamarByNomorKamar(c *gin.Context) {
	noStr := c.Param("num")
	var kamar models.KamarXTipeKamar

	query := `
		SELECT 
			k.id as id_kamar, k.nomor_kamar, tk.nama_tipe, k.status
			, k.created_at, k.updated_at
		FROM
			kamar k
		JOIN 
			tipe_kamar tk
		on k.id_tipe_kamar = tk.id
		WHERE
			k.nomor_kamar = $1
	`

	err := database.DBClient.Get(&kamar, query, noStr)

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

func GetKetersediaanKamarByDate(c *gin.Context) {
	tanggalMulai := c.Query("tanggal_mulai")
	tanggalSelesai := c.Query("tanggal_selesai")

	var kamar []models.KamarAvail

	query := `
	SELECT
	k.id as id_kamar, k.nomor_kamar, tk.nama_tipe, k.id_tipe_kamar,
	COALESCE(t_season.tarif, t_default.tarif) AS tarif,
	CASE
		WHEN s.id IS NOT NULL THEN s.id
		ELSE s_default.id
	END AS id_season
	FROM
	kamar k
	INNER JOIN
	tipe_kamar tk ON k.id_tipe_kamar = tk.id
	LEFT JOIN
	(
		SELECT * FROM tarif WHERE season_id IN (
			SELECT id FROM season WHERE $1 BETWEEN tanggal_mulai AND tanggal_berakhir
		)
	) t_season ON tk.id = t_season.id_tipe_kamar
	LEFT JOIN
	(
		SELECT * FROM tarif WHERE season_id IN (
			SELECT id FROM season WHERE nama_season = 'default'
		)
	) t_default ON tk.id = t_default.id_tipe_kamar
	LEFT JOIN
	season s ON t_season.season_id = s.id
	LEFT JOIN
	season s_default ON t_default.season_id = s_default.id
	LEFT JOIN
	reservasi r ON k.id= r.id_kamar
	WHERE
	k.status = true
	AND (
	(r.tanggal_checkin > $2 OR r.tanggal_checkout < $1 OR r.tanggal_checkin IS NULL OR r.tanggal_checkout IS NULL)
	)
`

	// SELECT
	// k.nomor_kamar, tk.nama_tipe,
	// COALESCE(t_season.tarif, t_default.tarif) AS tarif,
	// CASE
	// 	WHEN s.id IS NOT NULL THEN s.id
	// 	ELSE s_default.id
	// END AS id_season
	// FROM
	// kamar k
	// INNER JOIN
	// tipe_kamar tk ON k.id_tipe_kamar = tk.id
	// LEFT JOIN
	// (
	// 	SELECT * FROM tarif WHERE season_id IN (
	// 		SELECT id FROM season WHERE $1 BETWEEN tanggal_mulai AND tanggal_berakhir
	// 	)
	// ) t_season ON tk.id = t_season.id_tipe_kamar
	// LEFT JOIN
	// (
	// 	SELECT * FROM tarif WHERE season_id IN (
	// 		SELECT id FROM season WHERE nama_season = 'default'
	// 	)
	// ) t_default ON tk.id = t_default.id_tipe_kamar
	// LEFT JOIN
	// season s ON t_season.season_id = s.id
	// LEFT JOIN
	// season s_default ON t_default.season_id = s_default.id
	// LEFT JOIN
	// reservasi r ON k.nomor_kamar = r.nomor_kamar
	// WHERE
	// k.status = true
	// AND (
	// (r.tanggal_checkin > $2 OR r.tanggal_checkout < $1 OR r.tanggal_checkin IS NULL OR r.tanggal_checkout IS NULL)
	// )

	// SELECT
	//     k.nomor_kamar, tk.nama_tipe,  k.id as id_kamar, k.id_tipe_kamar,
	//     COALESCE(t_season.tarif, t_default.tarif) AS tarif
	// FROM
	//     kamar k
	// INNER JOIN
	//     tipe_kamar tk ON k.id_tipe_kamar = tk.id
	// LEFT JOIN
	//     (
	//         SELECT * FROM tarif WHERE season_id IN (
	//             SELECT id FROM season WHERE $1 BETWEEN tanggal_mulai AND tanggal_berakhir
	//         )
	//     ) t_season ON tk.id = t_season.id_tipe_kamar
	// LEFT JOIN
	//     (
	//         SELECT * FROM tarif WHERE season_id IN (
	//             SELECT id FROM season WHERE nama_season = 'default'
	//         )
	//     ) t_default ON tk.id = t_default.id_tipe_kamar
	// LEFT JOIN
	//     reservasi r ON k.nomor_kamar = r.nomor_kamar
	// WHERE
	//     k.status = true
	// AND (
	//     (r.tanggal_checkin > $2 OR r.tanggal_checkout < $1 OR r.tanggal_checkin IS NULL OR r.tanggal_checkout IS NULL)
	// )

	// SELECT
	// 		k.nomor_kamar, tk.nama_tipe, t.tarif
	//     FROM
	// 		kamar k
	//     INNER JOIN
	// 		tipe_kamar tk ON k.id_tipe_kamar = tk.id
	//     LEFT JOIN
	// 		tarif t ON tk.id = t.id_tipe_kamar
	// 	LEFT JOIN
	// 		reservasi r ON k.nomor_kamar = r.nomor_kamar
	//     WHERE
	// 		k.status = true
	//     AND (
	//         (r.tanggal_checkin > $1 AND r.tanggal_checkin > $2)
	//         OR (r.tanggal_checkout < $1 AND r.tanggal_checkout < $2)
	//         OR (r.tanggal_checkin IS NULL AND r.tanggal_checkout IS NULL)
	//     );

	err := database.DBClient.Select(&kamar, query, tanggalMulai, tanggalSelesai)

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
