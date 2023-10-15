package api

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

	r.GET("/api/user/:username", controllers.ShowUserDetailByIdParam)
	r.PUT("/api/user/update", middleware.Validate, controllers.UpdateProfile)

	r.GET("/api/tipekamar", middleware.Validate, controllers.GetTipeKamar)
	r.GET("/api/tipekamar/:id", middleware.Validate, controllers.GetTipeKamarById)
	r.POST("/api/tipekamar", middleware.Validate, controllers.CreateTipeKamar)
	r.PUT("/api/tipekamar/:id", middleware.Validate, controllers.UpdateTipeKamar)
	r.DELETE("/api/tipekamar/:id", middleware.Validate, controllers.DeleteTipeKamar)

	r.GET("/api/kamar", middleware.Validate, controllers.GetKamar)
	r.GET("/api/kamar/:id", middleware.Validate, controllers.GetKamarById)
	r.POST("/api/kamar", middleware.Validate, controllers.CreateKamar)
	r.PUT("/api/kamar/:id", middleware.Validate, controllers.UpdateKamar)
	r.DELETE("/api/kamar/:id", middleware.Validate, controllers.DeleteKamar)

	r.GET("/api/season", middleware.Validate, controllers.GetSeasons)
	r.GET("/api/season/:id", middleware.Validate, controllers.GetSeasonById)
	r.POST("/api/season", middleware.Validate, controllers.CreateSeason)
	r.PUT("/api/season/:id", middleware.Validate, controllers.UpdateSeason)
	r.DELETE("/api/season/:id", middleware.Validate, controllers.DeleteSeason)

	r.GET("/api/fasilitasberbayar", middleware.Validate, controllers.GetFasilitasBerbayars)
	r.GET("/api/fasilitasberbayar/:id", middleware.Validate, controllers.GetFasilitasBerbayarById)
	r.POST("/api/fasilitasberbayar", middleware.Validate, controllers.CreateFasilitasBerbayar)
	r.PUT("/api/fasilitasberbayar/:id", middleware.Validate, controllers.UpdateFasilitasBerbayar)
	r.DELETE("/api/fasilitasberbayar/:id", middleware.Validate, controllers.DeleteFasilitasBerbayar)

	r.GET("/api/tarif", middleware.Validate, controllers.GetTarifs)
	r.GET("/api/tarif/:id", middleware.Validate, controllers.GetTarifById)
	r.POST("/api/tarif", middleware.Validate, controllers.CreateTarif)
	r.PUT("/api/tarif/:id", middleware.Validate, controllers.UpdateTarif)
	r.DELETE("/api/tarif/:id", middleware.Validate, controllers.DeleteTarif)

	r.GET("/api/pong", middleware.Validate, controllers.ProtectedHandler)
}

// init gin app
func init() {
	app = gin.New()

	utility.InitToken("dnaidnaodnaw", 30)

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
