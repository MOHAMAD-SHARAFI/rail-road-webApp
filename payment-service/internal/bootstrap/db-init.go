package bootstrap

import (
	"log"
	"payment-service/internal/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	dsn := viper.GetString("dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "T_",
		}})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", models.ErrFailedConnectDB)
	}

	return db, nil
}

func Migration(db *gorm.DB) {
	err := db.AutoMigrate(&models.Payment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", models.ErrFailedMigrateDB)
	}
}
