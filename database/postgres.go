package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBClient *sqlx.DB

func ConnectPostgres() (db *sqlx.DB, err error) {

	// //onlinedb
	// host := "ep-jolly-waterfall-25791148.ap-southeast-1.aws.neon.tech"
	// port := "5432"
	// user := "GtFoBAE05"
	// pass := "VrcNhj62Fmbp"
	// dbname := "neondb"
	//local
	host := "localhost"
	port := "5432"
	user := "postgres"
	pass := "admin"
	dbname := "grandatma"
	// membuat data source name
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	// melakukan open koneksi ke postgres
	// driver postgres bisa di dapat dari melakukan import
	// dengan cara import _ "github.com/lib/pq"
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	// validate if db berhasil untuk terhubung
	// dengan cara melakukan `ping`
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DBClient = db
	return
}
