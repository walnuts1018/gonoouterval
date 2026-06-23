package sample

type Execer interface {
	Exec(query string) error
}

type DB struct{}

func (db *DB) Exec(query string) error {
	return nil
}

func RunTransaction(db Execer, fn func(tx Execer) error) error {
	return fn(db)
}

func SaveUser(db *DB) error {
	return RunTransaction(db, func(tx Execer) error {
		return db.Exec("INSERT INTO users (name) VALUES ('alice')") // want `using db from the outer scope but tx is defined inner at .*/transaction\.go:18`
	})
}
