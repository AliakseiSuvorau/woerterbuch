package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"woerterbuch-backend/src/model/requests"
	"woerterbuch-backend/src/modules/validation"
	"woerterbuch-backend/src/modules/web"
)

// RegisterWeb registers all routes for web.
func RegisterWeb(r *web.Router) {
	r.Group("", func() {
		r.Group("word", func() {
			r.Post("add", AddWord, validateRequest[requests.AddWordRequest])
			r.Post("edit", EditWord, validateRequest[requests.EditWordRequest])
			r.Delete("delete", DeleteWord, validateRequest[requests.DeleteWordRequest])
		})
		r.Group("dictionary", func() {
			r.Get("list", GetPage)
			r.Get("getRandomBatch", GetRandomBatch)
			r.Post("upload", Upload)
		})
	}, loggingMiddleware)
}

// loggingMiddleware logs start of a request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// validateRequest is used to validate requests.
func validateRequest[ReqT any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer bodyCloser(body)

		payload, readErr := io.ReadAll(body)
		if readErr != nil {
			log.Printf("Error has occurred while validating request body: %v", readErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var request ReqT
		if jsonErr := json.Unmarshal(payload, &request); jsonErr != nil {
			log.Printf("Error has occurred while validating request body: %v", jsonErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		validator := validation.CustomValidator{}
		if validationErr := validator.Validate(request); validationErr != nil {
			log.Printf("Error has occurred while validating request body: %v", validationErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(payload))
		next.ServeHTTP(w, r)
	})
}
