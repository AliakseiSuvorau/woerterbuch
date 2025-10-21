package repositories

import (
	"strconv"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model/dtos"
)

type WordsRepository struct{}

func (wr *WordsRepository) GetById(id uint64) (*dtos.Word, error) {
	var word dtos.Word
	result := global.DB.First(&word, strconv.FormatUint(id, 10))
	return &word, result.Error
}

func (wr *WordsRepository) GetByIds(ids []uint64) ([]*dtos.Word, error) {
	var words []*dtos.Word
	result := global.DB.Where("ID IN ?", ids).Find(&words)
	return words, result.Error
}

func (wr *WordsRepository) GetIds() ([]uint64, error) {
	var ids []uint64
	result := global.DB.Model(&dtos.Word{}).Pluck("id", &ids)
	return ids, result.Error
}

func (wr *WordsRepository) Insert(word *dtos.Word) error {
	if word == nil {
		return nil
	}
	result := global.DB.Create(word)
	return result.Error
}

func (wr *WordsRepository) Update(word *dtos.Word) error {
	if word == nil {
		return nil
	}
	result := global.DB.Save(word)
	return result.Error
}

func (wr *WordsRepository) GetAll() ([]*dtos.Word, error) {
	var words []*dtos.Word
	result := global.DB.Find(&words)
	return words, result.Error
}

// GetWordsByRange returns a list of words from 'from' inclusive to 'to' exclusive. Doesn't check limits.
func (wr *WordsRepository) GetWordsByRange(from, to int64) ([]*dtos.Word, error) {
	var words []*dtos.Word
	pageSize := to - from
	result := global.DB.Offset(int(from)).Limit(int(pageSize)).Find(&words)
	return words, result.Error
}

func (wr *WordsRepository) ExistsByWordAndTranslation(w *dtos.Word) bool {
	var count int64
	global.DB.Model(&dtos.Word{}).
		Where("word = ? AND translation = ?", w.Word, w.Translation).
		Count(&count)
	return count > 0
}

func (wr *WordsRepository) ExistsById(id uint64) bool {
	var count int64
	global.DB.Model(&dtos.Word{}).
		Where("id = ?", id).
		Count(&count)
	return count > 0
}

func (wr *WordsRepository) Count() int64 {
	var count int64
	global.DB.Table("words").Count(&count)
	return count
}

func (wr *WordsRepository) DeleteById(id uint64) error {
	result := global.DB.Delete(&dtos.Word{}, strconv.FormatUint(id, 10))
	return result.Error
}
