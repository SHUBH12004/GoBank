package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ShubhKanodia/GoBank/util"
	_ "github.com/lib/pq"
)

// dbDriver is the database driver used for testing.
// dbSource is the database source used for testing.

var testQueries *Queries

// testQueries is a global variable that holds the Queries object for testing.
// It is initialized in TestMain function which runs before any test cases.

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..") //parent folder
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// Load the database driver and source from the configuration.
	// Open a connection to the  database.

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	// Create a new Queries object .
	testQueries = New(testDB)
	os.Exit(m.Run())
}
