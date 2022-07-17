package vmilog

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

var (
	catlogMap     = make(map[string]*logrus.Logger)
	enableConsole bool
	enableFile    bool
	rwm           sync.RWMutex
	logPath       string
	level         = logrus.WarnLevel
)

func EnableConsole(enabled bool) {
	enableConsole = enabled
}

func EnableFile(enabled bool) {
	enableFile = enabled
}

func SetLogPath(p string) {
	logPath = p
}

func SetLevel(levelStr string) {
	switch levelStr {
	case INFO:
		level = logrus.InfoLevel
	case WARN:
		level = logrus.WarnLevel
	case ERROR:
		level = logrus.ErrorLevel
	default:
		level = logrus.WarnLevel
	}
}

func generateLogFileName(logCat string) string {
	if logPath == "" {
		return "logs/" + logCat + ".log"
	}
	return logPath + logCat + ".log"
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func getLogger(logCat string) *logrus.Logger {
	rwm.RLock()
	defer rwm.RUnlock()
	if l, finder := catlogMap[logCat]; finder {
		return l
	}
	return nil
}

func addLoggerIfNotExist(logCat string) *logrus.Logger {
	if lg := getLogger(logCat); lg != nil {
		return lg
	}
	rwm.Lock()
	defer rwm.Unlock()
	if l, finder := catlogMap[logCat]; finder {
		return l
	} else {
		log := logrus.New()
		var f *os.File
		if enableFile {
			fileName := generateLogFileName(logCat)
			dir := path.Dir(fileName)
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				os.MkdirAll(dir, 0700)
			}
			_f, err := os.Create(fileName)
			f = _f
			checkError(err)
		}
		if !enableConsole && !enableFile {
			log.SetOutput(io.Discard)
		}

		if enableConsole && !enableFile {
			log.SetOutput(os.Stdout)
		}
		if !enableConsole && enableFile {
			log.SetOutput(f)
		}
		if enableConsole && enableFile {
			mw := io.MultiWriter(os.Stdout, f)
			log.SetOutput(mw)
		}
		log.SetLevel(level)
		catlogMap[logCat] = log
		return log
	}
}

func Info(logCat string, v ...interface{}) {
	_log := addLoggerIfNotExist(logCat)
	_log.WithField("component", logCat).Info(v...)
}

func Warn(logCat string, v ...interface{}) {
	_log := addLoggerIfNotExist(logCat)
	_log.WithField("component", logCat).Warn(v...)
}

func Error(logCat string, v ...interface{}) {
	_log := addLoggerIfNotExist(logCat)
	_log.WithField("component", logCat).Error(v...)
}
