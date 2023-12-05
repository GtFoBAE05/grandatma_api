package main

import (
	"fmt"
	"grandatma_api/controllers"
	"grandatma_api/database"
	"grandatma_api/handler"
	"grandatma_api/middleware"
	"grandatma_api/utility"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func registerRouter(r *gin.RouterGroup) {
	r.GET("/api/ping", handler.Ping)

	r.POST("/api/auth/signup", controllers.CreatePengguna)
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/changepass", middleware.Validate, controllers.ChangePassword)

	r.GET("/api/user/:id", controllers.ShowUserDetailByIdParam)
	r.GET("/api/user/bytoken", middleware.Validate, controllers.ShowUserDetailByToken)
	r.GET("/api/user/customer/search", controllers.SearchCustomerByUsername)
	r.PUT("/api/user/update", middleware.Validate, controllers.UpdateProfile)

	r.GET("/api/tipekamar", middleware.Validate, controllers.GetTipeKamar)
	r.GET("/api/tipekamar/:id", middleware.Validate, controllers.GetTipeKamarById)
	r.POST("/api/tipekamar", middleware.Validate, controllers.CreateTipeKamar)
	r.PUT("/api/tipekamar/:id", middleware.Validate, controllers.UpdateTipeKamar)
	r.DELETE("/api/tipekamar/:id", middleware.Validate, controllers.DeleteTipeKamar)

	r.GET("/api/kamar", middleware.Validate, controllers.GetKamar)
	r.GET("/api/kamar/:id", middleware.Validate, controllers.GetKamarById)
	r.GET("/api/kamar/nomor/:num", middleware.Validate, controllers.GetKamarByNomorKamar)
	r.POST("/api/kamar", middleware.Validate, controllers.CreateKamar)
	r.PUT("/api/kamar/:id", middleware.Validate, controllers.UpdateKamar)
	r.DELETE("/api/kamar/:id", middleware.Validate, controllers.DeleteKamar)
	r.GET("/api/kamar/avail", middleware.Validate, controllers.GetKetersediaanKamarByDate)

	r.GET("/api/season", middleware.Validate, controllers.GetSeasons)
	r.GET("/api/season/:id", middleware.Validate, controllers.GetSeasonById)
	r.GET("/api/season/search", middleware.Validate, controllers.GetSeasonByName)
	r.POST("/api/season", middleware.Validate, controllers.CreateSeason)
	r.PUT("/api/season/:id", middleware.Validate, controllers.UpdateSeason)
	r.DELETE("/api/season/:id", middleware.Validate, controllers.DeleteSeason)

	r.GET("/api/fasilitasberbayar", middleware.Validate, controllers.GetFasilitasBerbayars)
	r.GET("/api/fasilitasberbayar/:id", middleware.Validate, controllers.GetFasilitasBerbayarById)
	r.GET("/api/fasilitasberbayar/search", middleware.Validate, controllers.GetFasilitasBerbayarByName)
	r.POST("/api/fasilitasberbayar", middleware.Validate, controllers.CreateFasilitasBerbayar)
	r.PUT("/api/fasilitasberbayar/:id", middleware.Validate, controllers.UpdateFasilitasBerbayar)
	r.DELETE("/api/fasilitasberbayar/:id", middleware.Validate, controllers.DeleteFasilitasBerbayar)

	r.POST("/api/fasilitasreservasi", middleware.Validate, controllers.CreateFasilitasReservasi)
	r.GET("/api/fasilitasreservasi/:id", middleware.Validate, controllers.GetFasilitasReservasiByIdReservasi)

	r.GET("/api/tarif", middleware.Validate, controllers.GetTarifs)
	r.GET("/api/tarif/:id", middleware.Validate, controllers.GetTarifById)
	r.POST("/api/tarif", middleware.Validate, controllers.CreateTarif)
	r.PUT("/api/tarif/:id", middleware.Validate, controllers.UpdateTarif)
	r.GET("/api/tarif/search", middleware.Validate, controllers.GetTarifByRoomTypeOrSeason)
	r.DELETE("/api/tarif/:id", middleware.Validate, controllers.DeleteTarif)

	r.POST("/api/reservasi", middleware.Validate, controllers.CreateReservasi)
	r.POST("/api/personal/reservasi", middleware.Validate, controllers.CreatePersonalReservasi)

	r.GET("/api/transaksihistory", middleware.Validate, controllers.GetTransaksiHistory)
	r.GET("/api/transaksihistory/:userId", middleware.Validate, controllers.GetTransaksiHistoryByUserId)
	r.GET("/api/transaksidetail/:id", middleware.Validate, controllers.GetTransaksiDetail)

	//jaminan
	r.PUT("/api/jaminan/:id", middleware.Validate, controllers.UpdateJaminan)
	r.PUT("/api/jaminanwithrek/:id", middleware.Validate, controllers.UpdateJaminanWithRekening)
	r.GET("/api/jaminan", middleware.Validate, controllers.GetJaminan)
	r.GET("/api/jaminan/:id", middleware.Validate, controllers.GetJaminanById)
	r.GET("/api/jaminan/my", middleware.Validate, controllers.GetMyJaminan)
	r.GET("/api/jaminan/my/:id", middleware.Validate, controllers.GetMyJaminanById)
	r.GET("/api/jaminan/search/uncompletejaminan", middleware.Validate, controllers.GetGroupUncompleteJaminanPayment)

	r.GET("/api/transaksi/search/batal", middleware.Validate, controllers.GetTransaksiByUsernameOrTransactionIdCanCancel)
	r.GET("/api/transaksi/my/search/batal", middleware.Validate, controllers.GetTransaksiByTransactionIdCanCancel)
	r.PUT("/api/transaksi/do/batalstatus/:id", middleware.Validate, controllers.UpdateStatusBatal)
	r.PUT("/api/transaksi/do/updatestatuslunas/:id", middleware.Validate, controllers.UpdateStatusLunas)

	//deposit
	r.PUT("/api/transaksi/do/updatedeposit/:id", middleware.Validate, controllers.UpdateStatusDeposit)

	//status menginap
	r.GET("/api/menginap/searchtransactionabletocheckin", middleware.Validate, controllers.SearchTransactionAbleToCheckin)
	r.GET("/api/menginap/searchtransactionabletocheckout", middleware.Validate, controllers.SearchTransactionAbleToCheckout)
	r.GET("/api/menginap/searchtransactionabletocomplete", middleware.Validate, controllers.SearchTransactionAbleToComplete)
	r.PUT("/api/menginap/do/updatecheckin/:id", middleware.Validate, controllers.UpdateStatusCheckin)
	r.PUT("/api/menginap/do/updatecheckout/:id", middleware.Validate, controllers.UpdateStatusCheckout)

	// uang muka grup
	// r.PUT("/api/transaksi/do/updateuangmuka/:id", middleware.Validate, controllers.UpdateUangMuka)

	//laporan
	r.GET("/api/report/newcustomerbyyear/:year", middleware.Validate, controllers.GetNewCustomerstatistics)
	r.GET("/api/report/monthlyreport/:year", middleware.Validate, controllers.GetMonthlyReport)
	r.GET("/api/report/visitorreport/:month/:year", middleware.Validate, controllers.GetVisitorStatistics)
	r.GET("/api/report/topcustomer/:year", middleware.Validate, controllers.GetTopCustomerByYear)

	r.GET("/api/pong", middleware.Validate, controllers.ProtectedHandler)
}

// init gin app
func init() {
	app = gin.Default()
	// app = gin.New()
	app.Use(CORSMiddleware())

	utility.InitToken("dnaidnaodnaw", 60)

	database.ConnectPostgres()

	// Handling routing errors
	app.NoRoute(func(c *gin.Context) {
		sb := &strings.Builder{}
		sb.WriteString("routing err: no route, try this:\n")
		for _, v := range app.Routes() {
			sb.WriteString(fmt.Sprintf("%s %s\n", v.Method, v.Path))
		}
		c.String(http.StatusBadRequest, sb.String())
	})

	r := app.Group("/")

	// register route
	registerRouter(r)
}

// entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// jika run menggunakan vercel dev, comment main
// func main() {
// 	app.Run(":53472")
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
