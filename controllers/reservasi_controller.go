package controllers

import (
	"fmt"
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateReservasi(c *gin.Context) {
	var reqBody models.Reservasi
	var totalReservasi models.TotalReservasi

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	reservasi := models.NewReservasi(reqBody.IdReservasi, reqBody.Email_pengguna, reqBody.IdKamar, reqBody.TanggalCheckin, reqBody.TanggalCheckout, reqBody.JumlahDewasa, reqBody.JumlahAnak, reqBody.NomorRekening, reqBody.PilihanKasur)

	//get total reservasi
	query := `
	SELECT 
		MAX(id) as total_reservasi
	FROM
		reservasi
	`

	err := database.DBClient.Get(&totalReservasi, query)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      true,
			"err disini": "select total reservasi",
			"message":    err.Error(),
		})
		return
	}

	jenisPelanggan := "P"
	if reqBody.IdReservasi == "grup" {
		jenisPelanggan = "G"
	}

	t := time.Now()
	tanggal := fmt.Sprintf("%d%02d%02d", t.Day(), t.Month(), t.Year()%100)

	totalReservasiStr := fmt.Sprintf("%d", totalReservasi.TotalReservasi+1)

	//gabungan id reservasi
	idReservasi := fmt.Sprintf("%s%s-%s", jenisPelanggan, tanggal, totalReservasiStr)

	//get id pengguna berdasarkan emai
	var users models.Pengguna
	query = `
		SELECT 
			id, email, password, role
			, created_at, updated_at
		FROM pengguna
		WHERE email = $1
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

	row := stmt.QueryRow(reqBody.Email_pengguna)

	err = row.Scan(
		&users.Id, &users.Email, &users.Password, &users.Role,
		&users.CreatedAt, &users.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": true,
			// "pesan":   "Tidak ditemukan email",
			"message": "Tidak ditemukan email",
		})
		return
	}

	//add reservasi
	query = `
	INSERT INTO reservasi (id_reservasi, id_pengguna, id_kamar, tanggal_checkin, tanggal_checkout, 
		jumlah_dewasa, jumlah_anak, nomor_rekening, pilihan_kasur, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

	_, err = stmt.Exec(idReservasi, users.Id, reservasi.IdKamar, reservasi.TanggalCheckin, reservasi.TanggalCheckout,
		reservasi.JumlahDewasa, reservasi.JumlahAnak, reservasi.NomorRekening, reservasi.PilihanKasur, reservasi.CreatedAt, reservasi.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get tarif kamar
	// var tarifKamar models.TarifKamar
	// query = `
	// 	SELECT
	// 		t.tarif as tarif
	// 	FROM
	// 		kamar k
	// 	JOIN
	// 		tarif t
	// 	ON
	// 		k.id_tipe_kamar = t.id_tipe_kamar
	// 	WHERE
	// 		k.id = $1
	// `

	// err = database.DBClient.Get(&tarifKamar, query, reservasi.Id)
	// if err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error":   true,
	// 		"message": err.Error(),
	// 		"error2":  "no rows",
	// 	})
	// 	return
	// }

	var kamar models.KamarAvail

	query = `
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
	reservasi r ON k.id = r.id_kamar
	WHERE
	k.status = true
	AND
	k.id = $2
	`

	err = database.DBClient.Get(&kamar, query, reservasi.TanggalCheckin, reservasi.IdKamar)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
			"error2":  "no rows",
		})
		return
	}

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
	totalPembayaran := float64(diffInDays) * kamar.Tarif

	//add deposit

	var idDeposit int64

	query = `
	INSERT INTO deposit (id_reservasi, nominal)
	VALUES ($1, $2)
	RETURNING id
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

	err = database.DBClient.QueryRow(query, idReservasi, 0).Scan(&idDeposit)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//add transaksi
	query = `
	INSERT INTO transaksi (id_reservasi,tanggal_transaksi, total_pembayaran, id_deposit, status_bayar, created_at, updated_at)
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

	_, err = stmt.Exec(idReservasi, time.Now(), totalPembayaran, idDeposit, false, reservasi.CreatedAt, reservasi.UpdatedAt)

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
		"data": gin.H{
			"id_reservasi": idReservasi,
		},
	})
}
