package driver

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DBDriver struct {
	Conn *gorm.DB
}

// Driver
func NewDBDriver() *DBDriver {
	DBMS := "postgres"
	USER := "postgres"
	PASS := ""
	DBNAME := ""
	HOST := "127.0.0.1"
	PORT := "5432"
	SSLMODE := "disable"

	DNS := "host=" + HOST + " port=" + PORT + " user=" + USER + " password=" + PASS + " dbname=" + DBNAME + " sslmode=" + SSLMODE

	var err error
	DB, err := gorm.Open(DBMS, DNS)
	if err != nil {
		panic(err)
	}

	db := new(DBDriver)
	db.Conn = DB
	return db
}
