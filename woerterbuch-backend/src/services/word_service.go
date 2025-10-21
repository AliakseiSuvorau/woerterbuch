package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model/dtos"
	"woerterbuch-backend/src/repositories"
)

type WordService struct {
	WordRepo repositories.WordsRepository
}

// AddWord checks if the word already exists and if not, adds it.
func (ws *WordService) AddWord(word *dtos.Word) error {
	if ws.WordRepo.ExistsByWordAndTranslation(word) {
		return nil
	}
	if err := ws.WordRepo.Insert(word); err != nil {
		return fmt.Errorf("add word: %w", err)
	}
	return nil
}

// GetWordsPage returns 'pageSize' words for page number 'pageNum'. Counts translates page number and page size into limits.
// Then checks the limits. If upper bound is greater than the total number of words, then returning words till the end.
// 'pageNum' is always greater than 0.
func (ws *WordService) GetWordsPage(pageNum, pageSize int64) ([]*dtos.Word, error) {
	from := (pageNum - 1) * pageSize // inclusive
	to := pageNum * pageSize         // exclusive

	var words []*dtos.Word

	// Check limits
	total := ws.WordRepo.Count()
	if from >= to || from > total {
		return words, nil
	}
	if to > total {
		to = total
	}

	words, err := ws.WordRepo.GetWordsByRange(from, to)
	if err != nil {
		return nil, fmt.Errorf("get '%d' words for page '%d' (from '%d' to '%d'): %w", pageSize, pageNum, from, to, err)
	}
	return words, nil
}

// GetRandomWords gets the batch size and checks if it is less than number of words in the dictionary,
// then retrieves 'batchSize' words from database and shuffles them.
func (ws *WordService) GetRandomWords() ([]*dtos.Word, error) {
	total := ws.WordRepo.Count()
	batchSize := min(total, global.WordRandomBatchSize)
	ids, err := ws.WordRepo.GetIds()
	if err != nil {
		return nil, fmt.Errorf("get ids: %w", err)
	}
	randomIndexes := ids[:batchSize]

	words, err := ws.WordRepo.GetByIds(randomIndexes)
	if err != nil {
		return nil, fmt.Errorf("get random words: %w", err)
	}

	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	return words, nil
}

// EditWord checks if the word exists by id and updates it.
func (ws *WordService) EditWord(word *dtos.Word) error {
	if !ws.WordRepo.ExistsById(word.ID) {
		return nil
	}

	if err := ws.WordRepo.Update(word); err != nil {
		return fmt.Errorf("edit word: %w", err)
	}
	return nil
}

// AddMultipleWords reads all words with articles and translations from csv reader, checks if they already exist
// and adds them if not.
func (ws *WordService) AddMultipleWords(reader *csv.Reader) error {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("add multiple words: %w", err)
		}

		newWord := dtos.Word{
			Article:     record[0],
			Word:        record[1],
			Translation: record[2],
		}
		if ws.WordRepo.ExistsByWordAndTranslation(&newWord) {
			continue
		}

		if insertErr := ws.WordRepo.Insert(&newWord); insertErr != nil {
			return fmt.Errorf("add multiple words: %w", insertErr)
		}
	}

	return nil
}

// DeleteById deletes word by id.
func (ws *WordService) DeleteById(id uint64) error {
	if err := ws.WordRepo.DeleteById(id); err != nil {
		return fmt.Errorf("delete word: %w", err)
	}
	return nil
}
