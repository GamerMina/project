package dbconn

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"projcet/types"
)

type AppConfig struct {
	Db types.DbConf
}

func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func LoadAppCfg() AppConfig {
	dsnstruct := types.DbConf{
		DbHost:     GetEnv("DB_HOST", "localhost"),
		DbUser:     GetEnv("DB_USER", "postgres"),
		DbPassword: GetEnv("DB_PASSWORD", "pass"),
		DbName:     GetEnv("DB_NAME", "postgres"),
		DbPort:     GetEnv("DB_PORT", "5432"),
	}
	return AppConfig{Db: dsnstruct}
}

func DbCon(dsnstruct types.DbConf) *gorm.DB {
	var dsn string
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dsnstruct.DbHost, dsnstruct.DbUser, dsnstruct.DbPassword, dsnstruct.DbName, dsnstruct.DbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	return db
}
