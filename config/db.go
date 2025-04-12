package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// _redis "github.com/go-redis/redis/v7"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type DB struct {
	*sql.DB
}

var db *gorm.DB

func DBInit() {

	dbinfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(dataSourceName string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	gormDB.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{postgres.Open(dataSourceName)},
			Policy:  dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(100).
			SetMaxOpenConns(150).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour),
	)

	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return gormDB, nil
}

func GetDB() *gorm.DB {
	return db
}
