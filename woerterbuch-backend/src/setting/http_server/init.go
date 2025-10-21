package http_server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/modules/web"
	"woerterbuch-backend/src/routes"
)

const defaultWordRandomBatchSize = 10

func Init() {
	r := web.NewRouter()
	routes.RegisterWeb(r)

	port := "6029" // Default port
	if len(os.Args) > 1 && os.Args[1] == "-port" {
		port = os.Args[2]
	}

	readRouteParameters()

	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func readRouteParameters() {
	wordRandomBatchSize, err := strconv.Atoi(os.Getenv("WORD_RANDOM_BATCH_SIZE"))
	if err != nil {
		global.WordRandomBatchSize = defaultWordRandomBatchSize
		global.Log.Printf("Error has occurred while retrieving word random batch size: %v", err)
		return
	}

	global.WordRandomBatchSize = int64(wordRandomBatchSize)
}
