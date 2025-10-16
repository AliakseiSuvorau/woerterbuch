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

func registerWebRoutes(r *web.Router) {
	addWordRoutes := func() {
		r.Group("word", func() {
			r.Post("add", routes.AddWord)
			r.Post("edit", routes.EditWord)
			r.Delete("delete", routes.DeleteWord)
		})
	}

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Incoming request:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	r.Group("dictionary", func() {
		addWordRoutes()
		r.Get("list", routes.GetPage)
		r.Get("getRandom", routes.GetRandomWords)
		r.Post("upload", routes.Upload)
	}, loggingMiddleware)
}

func Init() {
	r := web.NewRouter()
	registerWebRoutes(r)

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
