package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	host := "postgres"
	user := "postgres"
	password := "postgres"
	port := "5432"
	db_name := "cart_service_db"

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
