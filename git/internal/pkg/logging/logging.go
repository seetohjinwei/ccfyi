package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	stdOut   = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	errorOut = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
)

type MultiLogger struct{}

// Write should be unused.
func (l MultiLogger) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// WriteLevel uses the correct Writer based on the log level.
func (l MultiLogger) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return stdOut.Write(p)
	} else {
		return errorOut.Write(p)
	}
}

func Init() {
	setLogLevel()
	log.Logger = zerolog.New(MultiLogger{}).With().Timestamp().Logger()
}

const defaultLogLevel = "ERROR"

func setLogLevel() {
	logLevel := os.Getenv("LOG")
	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	switch logLevel {
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
}
