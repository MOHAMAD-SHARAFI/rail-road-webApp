package bootstrap

import (
	"user-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB(connectionString string) (*gorm.DB, error) {
	var db *gorm.DB
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "T_",
		}})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models.User{}, models.PassworResetToken{})
	if err != nil {
		return nil, err
	}

	return db, nil

}
