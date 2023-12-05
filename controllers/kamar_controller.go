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
	SELECT kamar.id as id_kamar, 
    kamar.id_tipe_kamar as id_tipe_kamar, 
    kamar.nomor_kamar as nomor_kamar, 
    tipe_kamar.nama_tipe as nama_tipe, 
    COALESCE(season_tarif.tarif, default_season.tarif) as tarif,
    COALESCE(season_tarif.season_id, default_season.season_id) as id_season,
	CASE
        WHEN transaksi.status_batal = false THEN 'Kamar tidak dapat dipakai'
        WHEN transaksi.status_batal = true THEN 'Kamar dapat dipakai'
        ELSE 'Unknown'
    END AS status_kamar
FROM kamar
LEFT JOIN tipe_kamar ON kamar.id_tipe_kamar = tipe_kamar.id
LEFT JOIN (
    SELECT DISTINCT ON (tarif.id_tipe_kamar) tarif.*, season.nama_season
    FROM tarif
    LEFT JOIN season ON tarif.season_id = season.id
    WHERE tarif.season_id = (
        SELECT id
        FROM season
        WHERE tanggal_mulai <= $1 
        AND tanggal_berakhir >= $1
    )
    ORDER BY tarif.id_tipe_kamar, season.tanggal_mulai DESC
) as season_tarif ON kamar.id_tipe_kamar = season_tarif.id_tipe_kamar
LEFT JOIN (
    SELECT tarif.*, season.nama_season
    FROM tarif
    LEFT JOIN season ON tarif.season_id = season.id
    WHERE season.nama_season = 'default'
) as default_season ON kamar.id_tipe_kamar = default_season.id_tipe_kamar
LEFT JOIN reservasi ON kamar.id = reservasi.id_kamar

LEFT JOIN transaksi ON reservasi.id_reservasi = transaksi.id_reservasi
WHERE kamar.id NOT IN (
    SELECT id_kamar
    FROM reservasi
    WHERE (tanggal_checkin <= $2 AND tanggal_checkout >= $1)
    OR (tanggal_checkin >= $1 AND tanggal_checkin <= $2)
    OR (tanggal_checkin <= $1 AND tanggal_checkout >= $2)
)
    OR transaksi.status_batal = true
GROUP BY kamar.id, kamar.id_tipe_kamar, kamar.nomor_kamar, tipe_kamar.nama_tipe, season_tarif.tarif, default_season.tarif, season_tarif.season_id, default_season.season_id, transaksi.status_batal;

`

	err := database.DBClient.Select(&kamar, query, tanggalMulai, tanggalSelesai)

	//tatau
	// 	SELECT DISTINCT
	//     kamar.id AS id_kamar,
	//     kamar.id_tipe_kamar AS id_tipe_kamar,
	//     kamar.nomor_kamar AS nomor_kamar,
	//     tipe_kamar.nama_tipe AS nama_tipe,
	//     CASE
	//         WHEN $1 >= season.tanggal_mulai AND $2 <= season.tanggal_berakhir THEN tarif.tarif
	//         ELSE (
	//             SELECT tarif
	//             FROM tarif
	//             INNER JOIN season ON tarif.season_id = season.id
	//             WHERE season.nama_season = 'default' AND tarif.id_tipe_kamar = kamar.id_tipe_kamar
	//         )
	//     END AS tarif,
	//     season.id AS id_season
	// FROM kamar
	// LEFT JOIN tipe_kamar ON kamar.id_tipe_kamar = tipe_kamar.id
	// LEFT JOIN tarif ON tarif.id_tipe_kamar = tipe_kamar.id
	// LEFT JOIN season ON tarif.season_id = season.id
	// LEFT JOIN reservasi ON reservasi.id_kamar = kamar.id
	// WHERE
	//     (
	//         $1 NOT BETWEEN reservasi.tanggal_checkin AND reservasi.tanggal_checkout
	//         OR $2 NOT BETWEEN reservasi.tanggal_checkin AND reservasi.tanggal_checkout
	//     )
	//     AND
	//     (
	//         reservasi.tanggal_checkin NOT BETWEEN $1 AND $2
	//         OR reservasi.tanggal_checkout NOT BETWEEN $1 AND $2
	//     )

	// SELECT
	// k.id, k.nomor_kamar, tk.nama_tipe, k.id_tipe_kamar,
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
	// reservasi r ON k.id = r.id_kamar
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
	// 		k.id, k.nomor_kamar, tk.nama_tipe, t.tarif
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
