package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

var DbPostgres *sqlx.DB

func PostgresConn() {
	hostEnv := os.Getenv("DB_HOST")
	portEnv := os.Getenv("DB_PORT")
	userEnv := os.Getenv("DB_USER")
	passEnv := os.Getenv("DB_PASS")
	dbnameEnv := os.Getenv("DB_NAME")

	sslMode := "require"
	caCertPath := "cerificate/ca-certificate.crt"

	DBStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s", hostEnv, portEnv, userEnv, passEnv, dbnameEnv, sslMode, caCertPath)
	db, err := sqlx.Connect("postgres", DBStr)

	CheckError(err)
	DbPostgres = db
	err = db.Ping()
	CheckError(err)
	fmt.Println("Postgres Connected!")
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error DB")
		panic(err)
	}
}
