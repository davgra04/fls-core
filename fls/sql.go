package fls

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // using sqlite
)

// DBAccessObject handles interactions with the database
type DBAccessObject struct {
	db *sqlx.DB
}

////////////////////////////////////////////////////////////////////////////////
// init
////////////////////////////////////////////////////////////////////////////////

func initSqliteDatabase(dbPath string) {
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

// NewSqliteDAO creates a new sqlite database access object and initializes the
// database itself if it does not exist
func NewSqliteDAO(dbPath string) *DBAccessObject {
	var err error

	if _, err := os.Stat(dbPath); err == nil {
		// path exists
		Info.Printf("using existing database at %v\n", dbPath)
	} else if os.IsNotExist(err) {
		// path does not exist
		Info.Printf("creating new database at %v\n", dbPath)
		initSqliteDatabase(dbPath)
	} else {
		// file may or may not exists, could be permissions?
		Error.Fatalf("problem determining if database at %v exists\n", dbPath)
	}

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("problem opening sqlite db: %v\n", err)
	}

	return &DBAccessObject{db: db}
}

////////////////////////////////////////////////////////////////////////////////
// DBAccessObject methods
////////////////////////////////////////////////////////////////////////////////

// Close closes the database handle
func (dao *DBAccessObject) Close() error {
	return dao.db.Close()
}

//////////////////////////////////////// artists table

func (dao *DBAccessObject) addArtist(artist *BandsInTownArtist) error {
	r, err := dao.db.NamedExec("INSERT INTO artists (name, bit_id, url, image_url, thumb_url, facebook_page_url, mbid, updated) VALUES (:name, :bit_id, :url, :image_url, :thumb_url, :facebook_page_url, :mbid, :updated)", artist)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	artist.ID = id
	return err
}

func (dao *DBAccessObject) getArtistList() ([]string, error) {
	artists := []string{}

	err := dao.db.Select(&artists, `SELECT name FROM artists`)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

//////////////////////////////////////// aliases table

func (dao *DBAccessObject) addAlias(alias, artistID string) error {
	_, err := dao.db.Exec("INSERT INTO aliases (alias, artist_id) VALUES (?, ?)", strings.ToLower(alias), artistID)
	return err
}

func (dao *DBAccessObject) getAliasList() ([]string, error) {
	aliases := []string{}

	err := dao.db.Select(&aliases, `SELECT alias FROM aliases`)
	if err != nil {
		return nil, err
	}

	return aliases, nil
}
