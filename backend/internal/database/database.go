package database

import (
	"database/sql"
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConnections struct {
	AppDB     *gorm.DB
	SamplesDB *sql.DB
}

func InitDBConnections(cfg *config.Config) (*DBConnections, error) {
	appDB, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize app database: %v", err)
	}

	samplesDB, err := initSamplesDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize samples database: %v", err)
	}

	return &DBConnections{
		AppDB:     appDB,
		SamplesDB: samplesDB,
	}, nil
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&models.User{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

func initSamplesDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&charset=utf8mb4&parseTime=True&loc=Local",
		cfg.SamplesDBUser,
		cfg.SamplesDBPassword,
		cfg.SamplesDBHost,
		cfg.SamplesDBPort,
		cfg.SamplesDBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to samples database: %v", err)
	}

	// test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping samples database: %v", err)
	}

	return db, nil
}
