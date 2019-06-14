package vault

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tablescript = "CREATE TABLE `encrypteditems` (`item` VARCHAR(64) PRIMARY KEY, `encrypted` BLOB NOT NULL);"
)

type VaultDatabase struct {
	path string
}

func Database(path string) *VaultDatabase {
	db := VaultDatabase{fmt.Sprintf("%s.db", path)}
	db.initialize()
	return &db
}

func (vault *VaultDatabase) initialize() error {
	path := vault.path
	// check db exists
	db, err := sql.Open("sqlite3", path)

	// check tables
	rows, err := db.Query("SELECT 1 FROM encrypteditems LIMIT 1;")
	if err != nil {
		// create tables
		_, err = db.Exec(tablescript)
		checkErr(err)
		log.Println("tables created")

	} else {
		rows.Close() //good habit to close
	}

	db.Close()

	return err
}

func (vault *VaultDatabase) transaction(callback func(db *sql.DB) error) error {
	db, err := sql.Open("sqlite3", vault.path)
	if err != nil {
		return err
	}

	defer db.Close()
	return callback(db)
}

func (vault *VaultDatabase) Item(item string) ([]byte, error) {

	var encrypted []byte
	err := vault.transaction(func(db *sql.DB) error {
		rows, e := db.Query("SELECT encrypted FROM encrypteditems WHERE item = ? LIMIT 1", item)
		checkErr(e)
		defer rows.Close() //good habit to close

		for rows.Next() {
			return rows.Scan(&encrypted)
		}

		return errors.New("Should have 1 row")
	})
	return encrypted, err
}

func (vault *VaultDatabase) SetItem(item string, encrypted []byte) error {
	return vault.transaction(func(db *sql.DB) error {
		_, err := db.Exec("INSERT INTO encrypteditems(item, encrypted) values(?,?)", item, encrypted)
		return err
	})
}

func (vault *VaultDatabase) Delete() {
	os.Remove(vault.path)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
