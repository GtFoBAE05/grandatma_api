package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetJaminan(c *gin.Context) {
	var jaminan []models.ListJaminan

	query := `
		SELECT 
			j.id, j.id_reservasi, j.nominal, j.status_lunas, t.total_pembayaran 
		FROM 
			jaminan j
		JOIN
			transaksi t
		ON
			j.id_reservasi = t.id_reservasi
		JOIN
			reservasi r
		ON 
			j.id_reservasi = r.id_reservasi
		WHERE
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
	`

	err := database.DBClient.Select(&jaminan, query)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,

			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  jaminan,
	})
}

func GetMyJaminan(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil ID pengguna.",
		})
		return
	}

	var jaminan []models.ListJaminan

	query := `
		SELECT 
			j.id, j.id_reservasi, j.nominal, j.status_lunas, t.total_pembayaran 
		FROM 
			jaminan j
		JOIN
			transaksi t
		ON
			j.id_reservasi = t.id_reservasi
		JOIN
			reservasi r
		ON
			j.id_reservasi = r.id_reservasi
		WHERE
			r.id_pengguna = $1
		AND
			r.tanggal_checkin > CURRENT_DATE
		AND
			t.status_batal = false
	`

	err := database.DBClient.Select(&jaminan, query, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,

			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  jaminan,
	})
}

func GetMyJaminanById(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil ID pengguna.",
		})
		return
	}

	idStr := c.Param("id")
	var jaminan models.ListJaminan

	query := `
		SELECT 
			j.id, j.id_reservasi, j.nominal, j.status_lunas, t.total_pembayaran 
		FROM 
			jaminan j
		JOIN
			transaksi t
		ON
			j.id_reservasi = t.id_reservasi
		JOIN
			reservasi r
		ON
			j.id_reservasi = r.id_reservasi
		WHERE 
			j.id_reservasi = $1 
		AND
			r.id_pengguna = $2
	`

	err := database.DBClient.Get(&jaminan, query, idStr, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,

			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  jaminan,
	})
}

func GetJaminanById(c *gin.Context) {

	idStr := c.Param("id")
	var jaminan models.ListJaminan

	query := `
		SELECT 
			j.id, j.id_reservasi, j.nominal, j.status_lunas, t.total_pembayaran 
		FROM 
			jaminan j
		JOIN
			transaksi t
		ON
			j.id_reservasi = t.id_reservasi
		WHERE j.id_reservasi = $1
	`

	err := database.DBClient.Get(&jaminan, query, idStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,

			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  jaminan,
	})
}

func GetGroupUncompleteJaminanPayment(c *gin.Context) {
	nama := c.Query("nama")
	id := c.Query("id")

	var transaksi []models.SearchUncompletedJaminan

	query := `
		SELECT 
			p.nama as nama, t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, 
			r.tanggal_checkin, j.nominal as total_jaminan, j.status_lunas as status_lunas
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
			jaminan j
		ON
			t.id_reservasi = j.id_reservasi
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

	var filteredTransaksi []models.SearchUncompletedJaminan

	for _, element := range transaksi {
		if strings.HasPrefix(element.IdReservasi, "G") {
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

func UpdateJaminan(c *gin.Context) {
	idStr := c.Param("id")
	var reqBody models.UpdateJaminanPayload
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get transaksi
	var transaksi models.TransaksiDetail

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin, r.tanggal_checkout,
			k.nomor_kamar, r.jumlah_dewasa, r.jumlah_anak, r.nomor_rekening, r.pilihan_kasur, t.status_batal
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
		WHERE
			t.id_reservasi = $1
	`
	err := database.DBClient.Get(&transaksi, query, idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	tipeCustomer := transaksi.IdReservasi[0]

	isFullPayment := false

	if tipeCustomer == 'P' {
		if reqBody.Nominal != transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":        true,
				"tipecustomer": "personal",
				"message":      "Nominal tidak sesuai",
			})
			return
		} else if reqBody.Nominal > transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Nominal tidak sesuai",
			})
			return
		} else {
			isFullPayment = true
		}
	} else {
		nominalMinimum := transaksi.TotalPembayaran * 0.5
		if reqBody.Nominal < nominalMinimum {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":        true,
				"tipecustomer": tipeCustomer,
				"grup":         "grup",
				"message":      "Nominal kurang dari 50% total pembayaran",
			})
			return
		}

		if reqBody.Nominal == transaksi.TotalPembayaran {
			isFullPayment = true
		} else if reqBody.Nominal > transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Nominal tidak sesuai",
			})
			return
		} else if reqBody.Nominal >= nominalMinimum {
			isFullPayment = false
		}
	}

	query = `
		UPDATE jaminan
		SET  nominal = $1, 
			status_lunas = $2,
			 updated_at = $3
		WHERE id_reservasi = $4
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

	_, err = stmt.Exec(reqBody.Nominal, isFullPayment, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated jaminan",
	})
}

func UpdateJaminanWithRekening(c *gin.Context) {
	idStr := c.Param("id")
	var reqBody models.UpdateJaminanWithRekeningPayload
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	//get transaksi
	var transaksi models.TransaksiDetail

	query := `
		SELECT 
			t.id_reservasi, t.total_pembayaran, t.tanggal_transaksi, r.tanggal_checkin, r.tanggal_checkout,
			k.nomor_kamar, r.jumlah_dewasa, r.jumlah_anak, r.nomor_rekening, r.pilihan_kasur, t.status_batal
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
		WHERE
			t.id_reservasi = $1
	`
	err := database.DBClient.Get(&transaksi, query, idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	tipeCustomer := transaksi.IdReservasi[0]

	isFullPayment := false

	if tipeCustomer == 'P' {
		if reqBody.Nominal != transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":        true,
				"tipecustomer": "personal",
				"message":      "Nominal tidak sesuai",
			})
			return
		} else if reqBody.Nominal > transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Nominal tidak sesuai",
			})
			return
		} else {
			isFullPayment = true
		}
	} else {
		nominalMinimum := transaksi.TotalPembayaran * 0.5
		if reqBody.Nominal < nominalMinimum {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":        true,
				"tipecustomer": tipeCustomer,
				"grup":         "grup",
				"message":      "Nominal kurang dari 50% total pembayaran",
			})
			return
		}

		if reqBody.Nominal == transaksi.TotalPembayaran {
			isFullPayment = true
		} else if reqBody.Nominal > transaksi.TotalPembayaran {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Nominal tidak sesuai",
			})
			return
		} else if reqBody.Nominal >= nominalMinimum {
			isFullPayment = false
		}
	}

	query = `
		UPDATE jaminan
		SET  nominal = $1, 
			status_lunas = $2,
			 updated_at = $3
		WHERE id_reservasi = $4
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

	_, err = stmt.Exec(reqBody.Nominal, isFullPayment, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query = `
		UPDATE reservasi
		SET  nomor_rekening = $1, 
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

	_, err = stmt.Exec(reqBody.Rekening, time.Now(), idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated jaminan",
	})
}
