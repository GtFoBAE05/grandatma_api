package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetNewCustomerstatistics(c *gin.Context) {

	year := c.Param("year")

	var customerCount []models.NewCustomerStatisticsByYear

	query := `
	SELECT
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 1) AS january,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 2) AS february,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 3) AS march,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 4) AS april,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 5) AS may,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 6) AS june,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 7) AS july,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 8) AS august,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 9) AS september,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 10) AS october,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 11) AS november,
    COUNT(*) FILTER (WHERE EXTRACT(MONTH FROM created_at) = 12) AS december,
	COUNT(*) AS total
FROM pengguna
WHERE EXTRACT(YEAR FROM created_at) = $1
AND role = 'customer'
	`

	err := database.DBClient.Select(&customerCount, query, year)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success get new customer statistics",
		"data":    customerCount,
	})
}

type MonthlyIncomeReport struct {
	Month int    `json:"month" db:"month"`
	Type  string `json:"type" db:"type"`
	Total int    `json:"total" db:"total"`
}

func GetMonthlyReport(c *gin.Context) {
	year := c.Param("year")

	var incomeReports []MonthlyIncomeReport
	result := make(map[string]map[string]int)
	response := []map[string]interface{}{}

	// List of month names
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	query := `
		SELECT
			EXTRACT(MONTH FROM r.tanggal_checkin) AS month,
			CASE
				WHEN LEFT(r.id_reservasi, 1) = 'P' THEN 'Personal'
				WHEN LEFT(r.id_reservasi, 1) = 'G' THEN 'Grup'
				ELSE 'Unknown'
			END AS type,
			SUM(t.total_pembayaran) AS total
		FROM transaksi t
		JOIN reservasi r ON t.id_reservasi = r.id_reservasi
		JOIN status_menginap sm ON t.id_reservasi = sm.id_reservasi
		WHERE EXTRACT(YEAR FROM r.tanggal_checkin) = $1 AND
		sm.status_checkin = true AND sm.status_checkout = true
		GROUP BY EXTRACT(MONTH FROM r.tanggal_checkin), type
		ORDER BY EXTRACT(MONTH FROM r.tanggal_checkin), type;
	`

	err := database.DBClient.Select(&incomeReports, query, year)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Initialize result map with zero values
	for _, month := range months {
		result[month] = map[string]int{"Personal": 0, "Grup": 0, "total": 0}
	}

	// Fill in values from the query result
	for _, report := range incomeReports {
		monthName := time.Month(report.Month).String()
		result[monthName][report.Type] = report.Total
		result[monthName]["total"] += report.Total
	}

	// Create array response
	for _, month := range months {
		monthData := make(map[string]interface{})
		monthData["month"] = month
		monthData["grup"] = result[month]["Grup"]
		monthData["personal"] = result[month]["Personal"]
		monthData["total"] = result[month]["total"]
		response = append(response, monthData)
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success get monthly income report",
		"data":    response,
	})
}

func GetVisitorStatistics(c *gin.Context) {
	month := c.Param("month")
	year := c.Param("year")

	var allRoomTypes []models.RoomType
	var roomCount []models.RoomVisitorReport

	// Query untuk mendapatkan semua tipe kamar yang mungkin
	allRoomTypesQuery := `
		SELECT id as id_tipe_kamar, nama_tipe as nama_tipe_kamar
		FROM tipe_kamar;
	`

	err := database.DBClient.Select(&allRoomTypes, allRoomTypesQuery)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Query utama untuk menghitung jumlah tamu berdasarkan jenis kamar
	visitorCountQuery := `
		SELECT
			tipe_kamar.nama_tipe AS type,
			COUNT(CASE WHEN LEFT(reservasi.id_reservasi, 1) = 'P' THEN 1 END) AS personal,
			COUNT(CASE WHEN LEFT(reservasi.id_reservasi, 1) = 'G' THEN 1 END) AS group,
			COUNT(*) AS total
		FROM
			tipe_kamar
		LEFT JOIN
			kamar ON tipe_kamar.id = kamar.id_tipe_kamar
		LEFT JOIN
			reservasi ON kamar.id = reservasi.id_kamar
		LEFT JOIN
			transaksi ON transaksi.id_reservasi = reservasi.id_reservasi
		WHERE
			EXTRACT(YEAR FROM reservasi.tanggal_checkin) = $1
			AND EXTRACT(MONTH FROM reservasi.tanggal_checkin) = $2
		GROUP BY
			tipe_kamar.nama_tipe;
	`

	err = database.DBClient.Select(&roomCount, visitorCountQuery, year, month)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Membuat peta untuk menyimpan hasil akhir
	result := make(map[string]models.RoomVisitorReport)

	// Inisialisasi nilai awal untuk setiap tipe kamar
	for _, roomType := range allRoomTypes {
		result[roomType.NamaTipeKamar] = models.RoomVisitorReport{
			Type:     roomType.NamaTipeKamar,
			Personal: 0,
			Group:    0,
			Total:    0,
		}
	}

	// Mengisi nilai sesuai dengan hasil query
	for _, room := range roomCount {
		result[room.Type] = room
	}

	// Mengonversi peta menjadi array
	var finalResult []models.RoomVisitorReport
	for _, v := range result {
		finalResult = append(finalResult, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success get room visitor statistics",
		"data":    finalResult,
	})
}

func GetTopCustomerByYear(c *gin.Context) {

	year := c.Param("year")

	var customerCount []models.TopCustomer

	query := `
	SELECT
    pengguna.nama as name,
    COUNT(reservasi.id_reservasi) AS reservation_count,
    SUM(transaksi.total_pembayaran) AS total_payment
FROM
    pengguna
JOIN
    reservasi ON pengguna.id = reservasi.id_pengguna
JOIN
    transaksi ON reservasi.id_reservasi = transaksi.id_reservasi
WHERE
	pengguna.role = 'customer'
AND
	EXTRACT(YEAR FROM transaksi.tanggal_transaksi) = $1
GROUP BY
    pengguna.nama
ORDER BY
    reservation_count DESC
LIMIT 5;
	`

	err := database.DBClient.Select(&customerCount, query, year)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success get new customer statistics",
		"data":    customerCount,
	})
}
