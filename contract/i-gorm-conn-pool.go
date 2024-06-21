package contract

import "gorm.io/gorm"

type IGromConnPool interface {
	ITraceable
	GetConn() (*gorm.DB, error)
}
