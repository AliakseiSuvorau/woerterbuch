package routes

import (
	"io"
	"log"
)

func bodyCloser(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Error has occurred while closing request body: %v", err)
	}
}
