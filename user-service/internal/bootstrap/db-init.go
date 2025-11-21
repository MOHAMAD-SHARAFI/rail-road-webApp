package bootstrap

import (
	"log"
	"user-service/internal/models"

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
		log.Fatalf("failed to connect database: %v", models.ErrFailedConnectDB)
	}
	return db, nil
}
func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", models.ErrFailedMigrateDB)
	}
}
