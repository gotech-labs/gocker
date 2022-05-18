package main

import (
	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/gocker/mysql"
)

func main() {
	db := mysql.New("5.7.38", "test_db")

	log.Info().
		Str("mysql", db.DatabaseURL()).
		Send()
}
