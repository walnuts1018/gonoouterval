package test1

import "test1/test1sub"

var _ = test1sub.Var

type Conn interface {
	Method()
}

type DB struct{}

func (db *DB) Method() {}

type Tx struct{}

func (tx *Tx) Method() {}

type Struct struct {
	DBField *DB
}

func runTransaction(func(*Tx)) {
}

func Foo() {
	var db *DB
	var c struct{ db *DB }

	runTransaction(func(tx *Tx) {
		tx.Method()
		db.Method()   // want `using db in the outer scope but tx is defined at .*/test1\.go:30`
		c.db.Method() // want `using c.db in the outer scope but tx is defined at .*/test1\.go:30`
		var _ Conn = nil
		var _ Struct = Struct{
			DBField: c.db, // want `using c.db in the outer scope but tx is defined at .*/test1\.go:30`
		}
	})
}
