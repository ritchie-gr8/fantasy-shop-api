package main

import (
	"github.com/ritchie-gr8/fantasy-shop-api/config"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/server"
)

func main() {
	conf := config.GetConfig()
	db := databases.NewPostgresDatabase(conf.Database)
	server := server.NewServer(conf, db)

	server.Start()
}
