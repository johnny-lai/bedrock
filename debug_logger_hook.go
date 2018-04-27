package bedrock

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type DebugLoggerHook struct {
}

func (hook *DebugLoggerHook) Fire(entry *log.Entry) error {
	_, file, line, ok := runtime.Caller(7)
	if ok == true {
		entry.Message = fmt.Sprintf("%s @ [%s:%d]", entry.Message, file, line)
	}

	return nil
}

func (hook *DebugLoggerHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
