package logger

import (
	"log"
	"os"

	"woerterbuch-backend/src/global"
)

func Init() {
	global.Log = log.New(os.Stdout, "[LOG]", log.Ldate|log.Ltime|log.Lshortfile)
}
