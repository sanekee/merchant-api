package log

import (
	"fmt"
	"log"
)

func Info(fmtStr string, args ...interface{}) {
	log.Printf("[INFO] %s\n", fmt.Sprintf(fmtStr, args...))
}

func Debug(fmtStr string, args ...interface{}) {
	log.Printf("[DEBUG] %s\n", fmt.Sprintf(fmtStr, args...))
}

func Error(fmtStr string, args ...interface{}) {
	log.Printf("[ERR] %s\n", fmt.Sprintf(fmtStr, args...))
}

func Warn(fmtStr string, args ...interface{}) {
	log.Printf("[ERR] %s\n", fmt.Sprintf(fmtStr, args...))
}
