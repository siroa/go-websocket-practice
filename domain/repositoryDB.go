package domain

import (
	"websocket/driver"
	"websocket/utils/errs"
)

type RepositoryDB struct {
	db driver.DBDriver
}

func (r RepositoryDB) saveMessage(msg Message) *errs.AppError {
	// ToDo: impl

	return nil
}

// factory function
func NewRepositoryDB(db *driver.DBDriver) RepositoryDB {
	return RepositoryDB{*db}
}

// auto migrate
func MyMigrate(db *driver.DBDriver) {
	//db.Conn.AutoMigrate()
}
