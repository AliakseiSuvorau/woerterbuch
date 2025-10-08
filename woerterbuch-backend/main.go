package main

import (
	"woerterbuch-backend/src/setting/database"
	"woerterbuch-backend/src/setting/http_server"
	"woerterbuch-backend/src/setting/logger"
)

func main() {
	logger.Init()
	database.Init()
	database.PrepareDictionary()
	http_server.Init()
}
