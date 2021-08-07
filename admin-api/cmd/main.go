package main

import (
	"fmt"
	"os"

	"admin-api/internal/api"

	_ "github.com/jackc/pgx/v4/stdlib"
	"xorm.io/xorm"
)

func main() {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("ADMIN_API_DB_USER"),
		os.Getenv("ADMIN_API_DB_PASSWORD"),
		os.Getenv("ADMIN_API_DB_HOST"),
		os.Getenv("ADMIN_API_DB_PORT"),
		os.Getenv("ADMIN_API_DB_NAME"),
	)
	fmt.Printf("Using DB Conn: '%s' for testing\n", conn)
	db, err := xorm.NewEngine("pgx", conn)
	if err != nil {
		panic(err)
	}

	api.Run(db)
}
