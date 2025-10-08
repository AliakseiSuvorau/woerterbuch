package global

import (
	"log"

	"gorm.io/gorm"
)

var Log *log.Logger

var DB *gorm.DB

var WordRandomBatchSize int64
