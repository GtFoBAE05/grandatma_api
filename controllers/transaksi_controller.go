package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"time"

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

func GetTransaksiByUsernameOrTransactionIdCanCancel(c *gin.Context) {
	nama := c.Query("nama")
	id := c.Query("id")

	var transaksi []models.SearchTransaksi

	query := `
		SELECT 
			p.nama, t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin
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
			p.nama LIKE '%' || $1 || '%'
		AND
			t.id_reservasi LIKE '%' || $2 || '%'
		AND
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
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

func GetTransaksiByTransactionIdCanCancel(c *gin.Context) {
	id := c.Query("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil ID pengguna.",
		})
		return
	}

	var transaksi []models.SearchTransaksi

	query := `
		SELECT 
			p.nama, t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin
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
			r.id_pengguna = $1
		AND
			t.id_reservasi LIKE '%' || $2 || '%'
		AND
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
	`
	err := database.DBClient.Select(&transaksi, query, userID, id)

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

// func GetTransaksiByUsernameOrTransactionIdNotCompletedPayment(c *gin.Context) {
// 	nama := c.Query("nama")
// 	id := c.Query("id")

// 	var transaksi []models.SearchUncompletedDeposit

// 	query := `
// 		SELECT
// 			p.nama, t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi,
// 			r.tanggal_checkin, d.nominal as total_deposit
// 		FROM
// 			transaksi t
// 		JOIN
// 			reservasi r
// 		ON
// 			t.id_reservasi = r.id_reservasi
// 		JOIN
// 			pengguna p
// 		ON
// 			r.id_pengguna = p.id
// 		JOIN
// 			deposit d
// 		ON
// 			t.id_reservasi = d.id_reservasi
// 		WHERE
// 			p.nama LIKE '%' || $1 || '%'
// 		AND
// 			t.id_reservasi LIKE '%' || $2 || '%'
// 		AND
// 			r.tanggal_checkin > CURRENT_DATE
// 		AND
// 			t.status_bayar = false
// 		AND
// 			t.status_batal = false
// 	`
// 	err := database.DBClient.Select(&transaksi, query, nama, id)

// 	if err != nil {
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
// 			"error":   true,
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	var filteredTransaksi []models.SearchUncompletedDeposit

// 	for _, element := range transaksi {
// 		if strings.HasPrefix(element.IdReservasi, "G") {
// 			filteredTransaksi = append(filteredTransaksi, element)
// 		}
// 	}

// 	if len(filteredTransaksi) == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error":   true,
// 			"message": "Data tidak ada",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"error": false,
// 		"data":  filteredTransaksi,
// 	})
// }

func GetTransaksiDetail(c *gin.Context) {
	transactionId := c.Param("id")

	var transaksi models.TransaksiDetail
	var fasilitasReservasi []models.FasilitasReservasiXTipeFasilitasXHarga
	var users models.Pengguna

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin, r.tanggal_checkout,
			k.nomor_kamar, r.jumlah_dewasa, r.jumlah_anak, r.nomor_rekening, r.pilihan_kasur, t.status_batal,
			k.id_tipe_kamar, tk.nama_tipe, j.nominal as jaminan, d.nominal as deposit
		FROM
			transaksi t
		JOIN
			reservasi r
		ON 
			t.id_reservasi = r.id_reservasi
		JOIN
			kamar k
		ON
			r.id_kamar = k.id
		JOIN
			tipe_kamar tk
		ON
			k.id_tipe_kamar = tk.id
		JOIN
			jaminan j
		ON 
			r.id_reservasi = j.id_reservasi
		JOIN
			deposit d
		ON 
			r.id_reservasi = d.id_reservasi
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

	//get fasilitas

	query = `
		SELECT 
			fr.id, fr.id_reservasi, fb.nama_fasilitas, fr.jumlah_unit, fb.harga , fb.created_at
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

	//get user
	query = `
		SELECT 
			p.nama, p.alamat
		FROM
			pengguna p
		JOIN
			reservasi r
		ON
			p.id = r.id_pengguna
		WHERE
			r.id_reservasi = $1
	`
	err = database.DBClient.Get(&users, query, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get tarif kamar normal

	tarif := transaksi.TotalPembayaran

	for _, element := range fasilitasReservasi {
		tarif -= float64(element.JumlahUnit) * float64(element.Harga)
	}

	//get selisih tanggal
	layout := "2006-01-02T15:04:05Z"
	t1, _ := time.Parse(layout, transaksi.TanggalCheckin)
	t2, _ := time.Parse(layout, transaksi.TanggalCheckout)

	difference := t2.Sub(t1)
	days := int(difference.Hours() / 24)

	c.JSON(http.StatusOK, gin.H{
		"error":     false,
		"data":      transaksi,
		"fasilitas": fasilitasReservasi,
		"user":      users,
		"tarif":     tarif,
		"days":      days,
	})
}

func UpdateStatusDeposit(c *gin.Context) {
	transactionId := c.Param("id")

	query := `
		UPDATE transaksi
		set  status_deposit = $1 
		WHERE id_reservasi = $2
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

	_, err = stmt.Exec(true, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated status deposit.",
	})
}

func UpdateUangMuka(c *gin.Context) {
	transactionId := c.Param("id")
	var reqBody models.UpdateDeposit

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE deposit
		set  nominal = $1 
		WHERE id_reservasi = $2
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

	_, err = stmt.Exec(reqBody.NominalDeposit, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated status bayar.",
	})
}

func UpdateStatusBatal(c *gin.Context) {
	transactionId := c.Param("id")

	query := `
		UPDATE transaksi
		set  status_batal = $1 
		WHERE id_reservasi = $2
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

	_, err = stmt.Exec(true, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated status batal.",
	})

}

func UpdateStatusLunas(c *gin.Context) {
	transactionId := c.Param("id")

	query := `
		UPDATE transaksi
		set  status_lunas = $1 
		WHERE id_reservasi = $2
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

	_, err = stmt.Exec(true, transactionId)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated lunas.",
	})
}
