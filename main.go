package main

import (
	"database/sql"

	"github.com/ShubhKanodia/GoBank/api"
	db "github.com/ShubhKanodia/GoBank/db/sqlc"
	"github.com/ShubhKanodia/GoBank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("Cannot load config: " + err.Error())
	}
	// This is the main function for the application.
	// It is currently empty because we are focusing on testing the database queries.
	// In a real application, you would typically start your server or run your application logic here.
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("Cannot con  nect to db: " + err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		panic("Cannot start server: " + err.Error())
	}
}
