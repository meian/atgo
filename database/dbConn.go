package database

import "gorm.io/gorm"

type DBConn struct {
	dbs []*gorm.DB
}

func NewDBConn(db *gorm.DB) *DBConn {
	return &DBConn{dbs: []*gorm.DB{db}}
}

func (r *DBConn) Transaction(f func(conn *DBConn) error) error {
	tx := r.dbs[0].Begin()
	if tx.Error != nil {
		return tx.Error
	}
	r.dbs = append(r.dbs, tx)
	defer func() {
		if r.dbs[len(r.dbs)-1] == tx {
			r.dbs = r.dbs[:len(r.dbs)-1]
		}
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := f(r); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *DBConn) DB() *gorm.DB {
	return r.dbs[len(r.dbs)-1]
}
