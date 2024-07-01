package initialisers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() Dbinstance {

	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger:      logger.Default.LogMode(logger.Info),
			PrepareStmt: true})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("db connected")
	db.Logger = logger.Default.LogMode(logger.Error)

	// log.Println("running migrations")
	// db.AutoMigrate(&models.Xuser{}, &models.UserWallet{}, &models.XuserAuthToken{}, &models.AdminUser{}, &models.AdminUserAuthToken{}, &models.Business{}, models.Event{}, models.EventTicket{})

	DB := Dbinstance{
		Db: db,
	}

	return DB

}
