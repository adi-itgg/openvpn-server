package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/goutil"
	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

// InitializeLogger initializes the zerolog logger.
func InitializeLogger(level zerolog.Level) {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i any) string {
			level := strings.ToUpper(goutil.String(i))
			switch level {
			case "DEBUG":
				return "\x1b[36mDEBUG\x1b[0m" // Cyan
			case "INFO":
				return "\x1b[37mINFO\x1b[0m" // White
			case "WARN":
				return "\x1b[33mWARN\x1b[0m" // Yellow
			case "ERROR":
				return "\x1b[31mERROR\x1b[0m" // Red
			case "FATAL":
				return "\x1b[35mFATAL\x1b[0m" // Magenta
			default:
				return level
			}
		},
		//FormatMessage: func(i any) string {
		//	return "|" + goutil.String(i) + "|"
		//},
		FormatCaller: func(i any) string {
			return filepath.Base(goutil.String(i))
		},
	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// This gets the module path (e.g., "github.com/your/module")
		// and helps in stripping the unnecessary prefix.
		// You might need to adjust this depending on your go.mod module path.
		_, _, _, ok := runtime.Caller(0)
		if !ok {
			return file // Fallback if we can't get caller info
		}

		// Find the last occurrence of the module path separator (e.g., "integrator-tfd/")
		// Or, a simpler approach is to find the last occurrence of "app/" or "shared-library/"
		// and then take the substring from there.

		// Let's try to get just the filename and line number
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1] + ":" + goutil.String(line)
		}
		return file + ":" + goutil.String(line)
	}

	log.Logger = zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Info().Msg("Logger initialized")
}
