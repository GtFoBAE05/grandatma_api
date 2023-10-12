package controllers

import (
	"grandatma_api/database"
	"grandatma_api/models"
	"grandatma_api/utility"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreatePengguna(c *gin.Context) {
	var reqBody models.Register

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	hashPass, err := utility.Hash(reqBody.Password)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	users := models.NewPengguna(reqBody.Nama, reqBody.Email, reqBody.Username, reqBody.Notelp, hashPass, reqBody.Role)

	query := `
	INSERT INTO pengguna (nama, email, username, notelp, password, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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

	_, err = stmt.Exec(users.Nama, users.Email, users.Username, users.Notelp, users.Password, users.Role, users.CreatedAt, users.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success create users",
	})

}

func Login(c *gin.Context) {
	var reqBody models.Login
	var users models.Pengguna

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
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

	row := stmt.QueryRow(reqBody.Email)

	err = row.Scan(
		&users.Id, &users.Email, &users.Password, &users.Role,
		&users.CreatedAt, &users.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"pesan":   "Tidak ditemukan email/pass",
			"message": err.Error(),
		})
		return
	}

	err = utility.Verify(users.Password, reqBody.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"pesan":   "Pass salah",
			"message": err.Error(),
		})
		return
	}

	token := utility.NewJWT(users.Id, users.Role)

	tokString, err := token.GenerateToken()

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"pesan":   "Pass salah",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": true,
		"token": tokString,
	})

}

func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Tidak ada token atau Gagal mengambil ID pengguna.",
		})
		return
	}
	var reqBody models.ChangePassword

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	hashPass, err := utility.Hash(reqBody.Password)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE pengguna
		set password = $1, updated_at = $2
		WHERE id = $3
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

	_, err = stmt.Exec(hashPass, time.Now(), userID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success updated password",
	})

}

func ShowUserDetailByIdParam(c *gin.Context) {
	username := c.Param("username")

	var user models.Pengguna

	query := `
		SELECT 
			id, nama, email
			, username, notelp
		FROM
			pengguna
		WHERE username = $1
	`
	err := database.DBClient.Get(&user, query, username)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  user,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Tidak ada token atau Gagal mengambil ID pengguna.",
		})
		return
	}
	var reqBody models.UpdateProfile

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	query := `
		UPDATE pengguna
		set nama = $1, 
		email = $2,
		username = $3,
		notelp = $4,
		updated_at = $5
		WHERE id = $6
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

	_, err = stmt.Exec(reqBody.Nama, reqBody.Email, reqBody.Notelp, reqBody.Notelp, time.Now(), userID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Success update profile",
	})
}

func ProtectedHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil ID pengguna.",
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil Role pengguna.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rute ini hanya dapat diakses jika token valid.",
		"user_id": userID,
		"role":    role,
	})
}
