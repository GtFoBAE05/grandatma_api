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

	users := models.NewPengguna(reqBody.Nama, reqBody.Email, reqBody.Username, reqBody.Notelp, hashPass, reqBody.Alamat, reqBody.Role)

	query := `
	INSERT INTO pengguna (nama, email, username, notelp, password, alamat,role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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

	_, err = stmt.Exec(users.Nama, users.Email, users.Username, users.Notelp, users.Password, users.Alamat, users.Role, users.CreatedAt, users.UpdatedAt)

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
			"pesan":   "Tidak ditemukan email",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"token": tokString,
		"role":  users.Role,
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

func ShowUserDetailByToken(c *gin.Context) {
	username, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Tidak ada token atau Gagal mengambil ID pengguna.",
		})
		return
	}

	var user models.Pengguna

	query := `
		SELECT 
			id, nama, email, alamat
			, username, notelp
		FROM
			pengguna
		WHERE id = $1
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

func ShowUserDetailByIdParam(c *gin.Context) {
	id := c.Param("id")

	var user models.Pengguna

	query := `
		SELECT 
			id, nama, email, alamat
			, username, notelp, role, created_at
		FROM
			pengguna
		WHERE id = $1
	`
	err := database.DBClient.Get(&user, query, id)

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

func SearchCustomerByUsername(c *gin.Context) {
	username := c.Query("username")

	var user []models.Pengguna

	query := `
		SELECT 
			id, nama, email, alamat
			, username, notelp, role
		FROM
			pengguna
		WHERE 
			username like '%' || $1 || '%'
		AND
			(role = 'customer' OR role = 'group')
	`
	err := database.DBClient.Select(&user, query, username)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if len(user) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Data tidak ada",
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
		alamat = $5,
		updated_at = $6
		WHERE id = $7
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

	_, err = stmt.Exec(reqBody.Nama, reqBody.Email, reqBody.Username, reqBody.Notelp, reqBody.Alamat, time.Now(), userID)

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
