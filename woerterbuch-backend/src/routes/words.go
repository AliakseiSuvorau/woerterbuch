package routes

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model"
	"woerterbuch-backend/src/repositories"
)

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
	if insertErr := wordRepository.Insert(&newWord); insertErr != nil {
		log.Printf("Error has occurred while inserting a new word: %v", insertErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GetPage(w http.ResponseWriter, r *http.Request) {
	pageNum, errConv := strconv.Atoi(r.URL.Query().Get("page"))
	if errConv != nil {
		log.Printf("Error has occurred while parsing page number: %v", errConv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	println(pageNum)

	pageSize, errConv := strconv.Atoi(r.URL.Query().Get("size"))
	if errConv != nil {
		log.Printf("Error has occurred while parsing page size: %v", errConv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	println(pageSize)

	wordRepository := repositories.WordsRepository{}
	words, err := wordRepository.GetRange((pageNum-1)*pageSize+1, pageNum*pageSize+1)
	if err != nil {
		log.Printf("Error has occurred while getting range of words: %v", err)
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

func castIntToUint64(arr []int) []uint64 {
	outArray := make([]uint64, 0, len(arr))
	for _, value := range arr {
		outArray = append(outArray, uint64(value))
	}
	return outArray
}

func processIndexes(arr []uint64) []uint64 {
	for i := range arr {
		arr[i]++
	}
	return arr
}

func GetRandomWords(w http.ResponseWriter, r *http.Request) {
	wr := repositories.WordsRepository{}
	numWordsInDict := wr.Count()
	batchSize := min(numWordsInDict, global.WordRandomBatchSize)

	randomGenerator := rand.New(rand.NewSource(time.Now().Unix()))
	randomIndexes := processIndexes(castIntToUint64(randomGenerator.Perm(int(numWordsInDict))[:batchSize]))

	words, err := wr.GetByIds(randomIndexes)
	if err != nil {
		log.Printf("Error has occurred while getting random words: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

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
	if insertErr := wordRepository.Update(&newWord); insertErr != nil {
		log.Printf("Error has occurred while inserting a new word: %v", insertErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	wordRepository := repositories.WordsRepository{}
	reader := csv.NewReader(body)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error has occurred while reading a record from csv: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newWord := model.Word{
			Article:     record[0],
			Word:        record[1],
			Translation: record[2],
		}
		if insertErr := wordRepository.Insert(&newWord); insertErr != nil {
			log.Printf("Error has occurred while inserting a new word: %v", insertErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func DeleteWord(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer bodyCloser(body)

	payload, readErr := io.ReadAll(body)
	if readErr != nil {
		log.Printf("Error has occurred while reading request body: %v", readErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordId := map[string]uint64{
		"id": 0,
	}
	if jsonErr := json.Unmarshal(payload, &wordId); jsonErr != nil {
		log.Printf("Error has occurred while unmarshalling request body: %v", jsonErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wordRepository := repositories.WordsRepository{}
	if deleteErr := wordRepository.DeleteById(wordId["id"]); deleteErr != nil {
		log.Printf("Error has occurred while deleting a word: %v", deleteErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
