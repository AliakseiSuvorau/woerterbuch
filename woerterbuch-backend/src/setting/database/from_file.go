package database

import (
	"encoding/csv"
	"os"
	"path"
	"strconv"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model"
	"woerterbuch-backend/src/repositories"
)

func PrepareDictionary() {
	readFromFile, errCastToBool := strconv.ParseBool(os.Getenv("READ_FROM_FILE"))
	if errCastToBool != nil || !readFromFile {
		global.Log.Print("Not using file for dictionary")
		return
	}
	dir := os.Getenv("DICT_DIR")
	filename := os.Getenv("DICT_FILENAME")
	filepath := path.Join(dir, filename)

	readDictionaryFromFile(filepath)
}

func readDictionaryFromFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		global.Log.Printf("Error has occurred while opening dictionary file: %v", err)
		return
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			global.Log.Printf("Error has occurred while closing dictionary file: %v", err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	wordsInfo, readErr := csvReader.ReadAll()
	if readErr != nil {
		global.Log.Printf("Error has occurred while reading dictionary file: %v", readErr)
		return
	}

	words := make([]*model.Word, 0, len(wordsInfo))
	for _, wordInfo := range wordsInfo {
		newWord := model.Word{
			Article:     wordInfo[0],
			Word:        wordInfo[1],
			Translation: wordInfo[2],
		}

		words = append(words, &newWord)
	}

	wr := repositories.WordsRepository{}
	if err := wr.InsertMultiple(words); err != nil {
		global.Log.Printf("Error has occurred while inserting words: %v", err)
		return
	}
}
