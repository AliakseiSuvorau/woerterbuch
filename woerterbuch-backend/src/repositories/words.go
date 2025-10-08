package repositories

import (
	"strconv"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model"
)

type WordsRepository struct{}

func (wr *WordsRepository) GetById(id uint64) (*model.Word, error) {
	var word model.Word
	result := global.DB.First(&word, strconv.FormatUint(id, 10))
	return &word, result.Error
}

func (wr *WordsRepository) GetByIds(ids []uint64) ([]*model.Word, error) {
	var words []*model.Word
	result := global.DB.Where("ID IN ?", ids).Find(&words)
	return words, result.Error
}

func (wr *WordsRepository) Insert(word *model.Word) error {
	if word == nil {
		return nil
	}

	if wr.existsByWordAndTranslation(word) {
		return nil
	}
	result := global.DB.Create(word)
	return result.Error
}

func (wr *WordsRepository) InsertMultiple(words []*model.Word) error {
	filteredWords := make([]*model.Word, 0, len(words))
	for _, word := range words {
		if word == nil {
			continue
		}

		if !wr.existsByWordAndTranslation(word) {
			filteredWords = append(filteredWords, word)
		}
	}

	if len(filteredWords) == 0 {
		return nil
	}

	result := global.DB.Create(filteredWords)
	return result.Error
}

func (wr *WordsRepository) Update(word *model.Word) error {
	if word == nil {
		return nil
	}

	if !wr.existsById(word) {
		return nil
	}

	result := global.DB.Save(word)
	return result.Error
}

func (wr *WordsRepository) GetAll() ([]*model.Word, error) {
	var words []*model.Word
	result := global.DB.Find(&words)
	return words, result.Error
}

// GetRange returns a list of words. If from = 1 and to = 8, it will return words from 0 to
func (wr *WordsRepository) GetRange(from, to int) ([]*model.Word, error) {
	var words []*model.Word
	total := int(wr.Count())
	if from > to || from > total {
		return words, nil
	}

	if to > total {
		to = total
	}

	from--
	pageSize := to - from
	result := global.DB.Offset(from).Limit(pageSize).Find(&words)
	return words, result.Error
}

func (wr *WordsRepository) existsByWordAndTranslation(w *model.Word) bool {
	var count int64
	global.DB.Model(&model.Word{}).
		Where("word = ? AND translation = ?", w.Word, w.Translation).
		Count(&count)
	return count > 0
}

func (wr *WordsRepository) existsById(w *model.Word) bool {
	var count int64
	global.DB.Model(&model.Word{}).
		Where("id = ?", w.ID).
		Count(&count)
	return count > 0
}

func (wr *WordsRepository) Count() int64 {
	var count int64
	global.DB.Table("words").Count(&count)
	return count
}

func (wr *WordsRepository) DeleteById(id uint64) error {
	result := global.DB.Delete(&model.Word{}, strconv.FormatUint(id, 10))
	return result.Error
}
