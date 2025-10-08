package main

import (
	"woerterbuch-backend/src/setting/database"
	"woerterbuch-backend/src/setting/http_server"
	"woerterbuch-backend/src/setting/logger"
)

func main() {
	logger.Init()
	database.Init()
	http_server.Init()
}
