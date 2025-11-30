package bootstrap

import (
	"user-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "T_",
		}})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models.User{}, models.PasswordResetToken{})
	if err != nil {
		return nil, err
	}

	return db, nil

}
