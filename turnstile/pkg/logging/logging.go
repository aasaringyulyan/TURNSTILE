package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

type Logger struct {
	*logrus.Entry
}

func (s *Logger) ExtraFields(fields map[string]interface{}) *Logger {
	return &Logger{s.WithFields(fields)}
}

var instance Logger
var once sync.Once

func GetLogger(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}

		l := logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
			},
			DisableColors:   false,
			TimestampFormat: time.Stamp,
			FullTimestamp:   true,
		}

		if logrusLevel != logrus.DebugLevel {
			err = os.MkdirAll("logs", 0700)
			if err != nil {
				panic(err)
			}

			allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
			if err != nil {
				fmt.Printf("err: %s", err)
			}

			l.AddHook(&writerHook{
				Writer:    []io.Writer{allFile},
				LogLevels: logrus.AllLevels,
			})
		}

		l.SetOutput(os.Stdout)
		l.SetLevel(logrusLevel)

		instance = Logger{logrus.NewEntry(l)}
	})

	return instance
}
