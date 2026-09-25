package log

import (
	"fmt"
	"os"
	"strings"
)

const envVarName = "DEBUG"

var GlobalLogger = New()

type Logger struct {
	enabled bool
}

func isDebugEnabled() bool {
	val, ok := os.LookupEnv(envVarName)
	if ok {
		val = strings.TrimSpace(val)
		val = strings.ToLower(val)

		if val == "" || val == "0" || val == "false" {
			return false
		}

		return true
	}

	return false
}

func New() *Logger {
	return &Logger{
		enabled: isDebugEnabled(),
	}
}

func (l *Logger) Log(str string) {
	if l.enabled {
		fmt.Println(str)
	}
}

func (l *Logger) Logf(format string, a ...any) {
	if l.enabled {
		fmt.Printf(format, a...)
	}
}

func Log(str string) {
	GlobalLogger.Log(str)
}

func Logf(format string, a ...any) {
	GlobalLogger.Logf(format, a...)
}
