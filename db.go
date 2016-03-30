package db

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db   *gorm.DB
	once sync.Once
)

func connect() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Panic("Failed to open database connection to postgres: ", err)
	}

	err = db.DB().Ping()

	if err != nil {
		log.Panic("Dead connection to postgres: ", err)
	}

	if maxIdleConnections := os.Getenv("DATABASE_MAX_IDLE_CONNECTIONS"); maxIdleConnections == "" {
		db.DB().SetMaxIdleConns(1)
	} else if mc, err := strconv.Atoi(maxIdleConnections); err == nil {
		db.DB().SetMaxIdleConns(mc)
	} else {
		log.Panic("Could not convert DATABASE_MAX_IDLE_CONNECTIONS to int: ", err)
	}
	if maxConnections := os.Getenv("DATABASE_MAX_CONNECTIONS"); maxConnections == "" {
		db.DB().SetMaxOpenConns(100)
	} else if mc, err := strconv.Atoi(maxConnections); err == nil {
		db.DB().SetMaxOpenConns(mc)
	} else {
		log.Panic("Could not convert DATABASE_MAX_CONNECTIONS to int: ", err)
	}
}

func DB() *gorm.DB {
	once.Do(connect)
	return db
}
