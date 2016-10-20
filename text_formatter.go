// Based on vendor/github.com/Sirupsen/logrus/text_formatter.go
package bedrock

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"sort"
	"time"
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type LogTextFormatter struct {
	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool
}

func (f *LogTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	b := &bytes.Buffer{}

	// prefixFieldClashes(entry.Data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = log.DefaultTimestampFormat
	}

	f.printEntry(b, entry, keys, timestampFormat)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *LogTextFormatter) printEntry(b *bytes.Buffer, entry *log.Entry, keys []string, timestampFormat string) {
	if !f.FullTimestamp {
		fmt.Fprintf(b, "[%04d] %-44s ", miniTS(), entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] %-44s ", entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s=%+v", k, v)
	}
}
