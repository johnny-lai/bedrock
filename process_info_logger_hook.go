package bedrock

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

type ProcessInfoLoggerHook struct {
}

func (hook *ProcessInfoLoggerHook) Fire(entry *log.Entry) error {
	entry.Message = fmt.Sprintf("(pid=%d) %s", os.Getpid(), entry.Message)

	return nil
}

func (hook *ProcessInfoLoggerHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
