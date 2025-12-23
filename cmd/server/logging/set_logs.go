package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogging() {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	consoleWriter.FormatTimestamp = func(i interface{}) string {
		t, err := time.Parse(time.RFC3339Nano, i.(string))
		if err != nil {
			return ""
		}
		return t.In(loc).Format("2006-01-02 15:04:05")
	}

	if os.Getenv("ENV") == "prod" {
		log.Logger = zerolog.New(consoleWriter).
			With().
			Timestamp().
			Logger()
	} else {
		log.Logger = zerolog.New(consoleWriter).
			With().
			Timestamp().
			Caller(). // agrega archivo y línea (opcional)
			Logger()
	}

	// Nivel mínimo de log
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// package logging

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"
// )

// // Colores ANSI
// const (
// 	reset      = "\033[0m"
// 	red        = "\033[31m"
// 	green      = "\033[32m"
// 	yellow     = "\033[33m"
// 	cyan       = "\033[36m"
// 	whiteBgRed = "\033[41;97m"
// )

// var logger = log.New(os.Stdout, "", 0)

// func logMessage(level string, color string, format string, args ...any) {
// 	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
// 	timestamp := time.Now().In(loc).Format("2006-01-02 15:04:05")
// 	message := fmt.Sprintf(format, args...)
// 	formatted := fmt.Sprintf("%s%s - %s - %s%s", color, timestamp, level, message, reset)
// 	logger.Println(formatted)
// }

// func DEBUG(format string, args ...any) {
// 	logMessage("DEBUG", cyan, format, args...)
// }

// func INFO(format string, args ...any) {
// 	logMessage("INFO", green, format, args...)
// }

// func WARNING(format string, args ...any) {
// 	logMessage("WARNING", yellow, format, args...)
// }

// func ERROR(format string, args ...any) {
// 	logMessage("ERROR", red, format, args...)
// }

// func CRITICAL(format string, args ...any) {
// 	logMessage("CRITICAL", whiteBgRed, format, args...)
// }
