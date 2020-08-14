package fls

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3" // using sqlite
)

// DB is the global sql database object
var DB *sql.DB
var dbPath string = "data.sql"

////////////////////////////////////////////////////////////////////////////////
// DB access methods
////////////////////////////////////////////////////////////////////////////////

func addArtists(artists []string) error {
	if len(artists) == 0 {
		return nil
	}

	// being transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// perform inserts
	for _, a := range artists {

		a = strings.TrimSpace(a)

		rows, err := tx.Query(`SELECT name FROM artists WHERE LOWER(name) LIKE ?`, strings.ToLower(a))
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if rows.Next() {
			Info.Printf("    already have artist entry for %v\n", a)
			continue
		}

		_, err = tx.Exec(`INSERT INTO artists (name) VALUES (?)`, a)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// complete transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getArtists() ([]string, error) {
	rows, err := DB.Query(`SELECT name FROM artists`)
	if err != nil {
		return nil, err
	}

	artists := []string{}

	for rows.Next() {
		var a string
		rows.Scan(&a)
		artists = append(artists, a)
	}

	return artists, nil
}

////////////////////////////////////////////////////////////////////////////////
// init
////////////////////////////////////////////////////////////////////////////////

func createNewDatabase() {
	schemaPath := "sql/schema0.sql"
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("problem reading schema sql: %v\n", err)
	}

	// dsn := fmt.Sprintf("file:%v", *dbPath)
	DB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("problem opening sqlite db: %v\n", err)
	}
	defer DB.Close()

	if _, err = DB.Exec(string(schemaSQL)); err != nil {
		log.Fatalf("problem initializing db: %v\n", err)
	}
}

// InitializeDatabase creates a new sqlite database if none exists
func InitializeDatabase() {
	var err error

	if _, err := os.Stat(dbPath); err == nil {
		// path exists
		Info.Printf("using existing database at %v\n", dbPath)
	} else if os.IsNotExist(err) {
		// path does not exist
		Info.Printf("creating new database at %v\n", dbPath)
		createNewDatabase()
	} else {
		// file may or may not exists, could be permissions?
		Error.Fatalf("problem determining if database at %v exists\n", dbPath)
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("problem opening sqlite db: %v\n", err)
	}
}
