package db

import (
	"github.com/vitorpcruz/desafios-golang/client-server-api/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConfigureDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./coins.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Coin{})

	return db
}
