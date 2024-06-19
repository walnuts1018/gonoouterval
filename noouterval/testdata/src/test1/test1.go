package test1

import "test1/test1sub"

type DB struct{}

func (db *DB) Method() {}

type Tx struct{}

func (tx *Tx) Method() {}

type Struct struct {
	DBField *DB
}

type DeepStruct struct {
	s Struct
}

func runTransaction(func(*Tx)) {
}

func Foo() {
	var db *DB
	var c struct{ db *DB }
	var d DeepStruct

	runTransaction(func(tx *Tx) {
		tx.Method()
		db.Method()   // want `using db from the outer scope but tx is defined inner at .*/test1\.go:29`
		c.db.Method() // want `using c.db from the outer scope but tx is defined inner at .*/test1\.go:29`
		var _ test1sub.Conn = nil
		var _ Struct = Struct{
			DBField: c.db, // want `using c.db from the outer scope but tx is defined inner at .*/test1\.go:29`
		}
		d.s.DBField.Method() // want `using d.s.DBField from the outer scope but tx is defined inner at .*/test1\.go:29`
	})
}
