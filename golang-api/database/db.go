package database

import (
	"fmt"
	"log"

	"github.com/fixia/golang-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect abre la conexión a MySQL. El esquema es gestionado por Flyway
// (ver ./migrations y el servicio `flyway` en docker-compose.yml).
// No se ejecuta AutoMigrate — las tablas deben existir antes de arrancar la app.
func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}

	DB = db
	log.Println("✅ Database connected")
	return db
}
