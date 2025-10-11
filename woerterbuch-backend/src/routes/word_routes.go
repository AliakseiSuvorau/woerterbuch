package routes

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"woerterbuch-backend/src/model"
	"woerterbuch-backend/src/repositories"
	"woerterbuch-backend/src/services"
)

// AddWord handles add word request. Parses request and adds a word with specified article and translation.
func AddWord(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	payload, readErr := io.ReadAll(body)
	if readErr != nil {
		log.Printf("Error has occurred while reading request body: %v", readErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newWord model.Word
	if jsonErr := json.Unmarshal(payload, &newWord); jsonErr != nil {
		log.Printf("Error has occurred while unmarshalling request body: %v", jsonErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	if insertErr := wordService.AddWord(&newWord); insertErr != nil {
		log.Printf("Error has occurred while inserting a new word: %v", insertErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetPage handles request of getting page of words. Parses request and returns 'pageSize' words for page
// with number 'pageNum'. Then wraps data and sends response back.
func GetPage(w http.ResponseWriter, r *http.Request) {
	pageNum, errConv := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if errConv != nil {
		log.Printf("Error has occurred while parsing page number: %v", errConv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageSize, errConv := strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
	if errConv != nil {
		log.Printf("Error has occurred while parsing page size: %v", errConv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	words, err := wordService.GetWordsPage(pageNum, pageSize)
	if err != nil {
		log.Printf("Error has occurred while getting page of words: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, jsonMarshalErr := json.Marshal(words)
	if jsonMarshalErr != nil {
		log.Printf("Error has occurred while marshalling response: %v", jsonMarshalErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, responseWriteErr := w.Write(response); responseWriteErr != nil {
		log.Printf("Error has occurred while writing response: %v", responseWriteErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// GetRandomWords handles request for getting random words. Retrieves a random batch of a fixed size of
// words from dictionary. Wraps and sends response.
func GetRandomWords(w http.ResponseWriter, r *http.Request) {
	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	words, err := wordService.GetRandomWords()
	if err != nil {
		log.Printf("Error has occurred while getting random words: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, jsonMarshalErr := json.Marshal(words)
	if jsonMarshalErr != nil {
		log.Printf("Error has occurred while marshalling response: %v", jsonMarshalErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, responseWriteErr := w.Write(response); responseWriteErr != nil {
		log.Printf("Error has occurred while writing response: %v", responseWriteErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// EditWord handles edit word request. Parses request and edits the word in dictionary.
func EditWord(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	payload, readErr := io.ReadAll(body)
	if readErr != nil {
		log.Printf("Error has occurred while reading request body: %v", readErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newWord model.Word
	if jsonErr := json.Unmarshal(payload, &newWord); jsonErr != nil {
		log.Printf("Error has occurred while unmarshalling request body: %v", jsonErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	if editErr := wordService.EditWord(&newWord); editErr != nil {
		log.Printf("Error has occurred while inserting a new word: %v", editErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Upload handles dictionary csv-file upload request.
func Upload(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	if err := wordService.AddMultipleWords(csv.NewReader(body)); err != nil {
		log.Printf("Error has occurred while uploading dictionary: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// DeleteWord handles delete word request.
func DeleteWord(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	payload, readErr := io.ReadAll(body)
	if readErr != nil {
		log.Printf("Error has occurred while reading request body: %v", readErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordId := map[string]uint64{"id": 0}
	if jsonErr := json.Unmarshal(payload, &wordId); jsonErr != nil {
		log.Printf("Error has occurred while unmarshalling request body: %v", jsonErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordRepository := repositories.WordsRepository{}
	wordService := services.WordService{WordRepo: wordRepository}
	if deleteErr := wordService.DeleteById(wordId["id"]); deleteErr != nil {
		log.Printf("Error has occurred while deleting a word: %v", deleteErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
