package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

/*************************************************
Debug Level Setting
- debug
- info
- warning
- error
- fatal
- panic
*************************************************/

const (
	fileTag = "file"
	lineTag = "line"
	funcTag = "func"
)

type logConf struct { //執行時期的 log 功能配置
	showFileInfo bool //是否顯示 file name, func name, line number
	today        *time.Time
	logpath      string
	file         *os.File
}

var rtLogConf logConf
var logExit = make(chan error)
var wg sync.WaitGroup

// Init ...
func Init(env, level, logpath, duration, url, channel string, hook, forceColor, fullTimestamp bool) {
	lv, _ := logrus.ParseLevel(level)

	format := &logrus.TextFormatter{ForceColors: forceColor, FullTimestamp: fullTimestamp}

	switch env {
	case "dev":
		InitLog(format, lv, env, logpath, duration, url, channel, hook, true, true)
	case "uat":
		InitLog(format, lv, env, logpath, duration, url, channel, hook, false, true)
	case "prod":
		InitLog(format, lv, env, logpath, duration, url, channel, hook, false, true)
	}
}

//InitLog config the log
func InitLog(format logrus.Formatter, level logrus.Level, env, logpath, duration, url, channel string, hook, multiWriter, showFileInfo bool) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(fmt.Sprintf("InitLog %v", err))
	}

	logrus.SetFormatter(format)
	logrus.SetLevel(level)

	if hook {
		// logrus.AddHook(&slack.Hook{
		// 	HookURL:        url,
		// 	AcceptedLevels: slack.LevelThreshold(level),
		// 	Channel:        channel,
		// 	IconEmoji:      ":ghost:",
		// 	Username:       "footbot",
		// 	Env:            env,
		// })
	}

	now := time.Now().UTC()
	t := now.Truncate(d)

	fullpath := logpath + "/" + Time2String(&t, "_") + "." + "log"

	if err := os.MkdirAll(filepath.Dir(fullpath), 0744); err != nil {
		Error("error folder create : ", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(fullpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		Error("error opening file: ", err)
		os.Exit(1)
	}

	if multiWriter {
		logrus.SetOutput(io.MultiWriter(f, os.Stdout))
	} else {
		logrus.SetOutput(f)
	}

	rtLogConf.showFileInfo = showFileInfo
	rtLogConf.today = &t
	rtLogConf.logpath = logpath
	rtLogConf.file = f

	wg.Add(1)
	go func() {
		Debug("Logger start fetching filename with datetime")
	loop:
		for {
			now := time.Now().UTC()
			next := now.Truncate(d)
			nt := next.Unix()
			t := rtLogConf.today.UTC().Unix()

			if nt > t {
				path := rtLogConf.logpath + "/" + Time2String(&next, "_") + "." + "log"
				f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					panic(err)
				}

				if multiWriter {
					logrus.SetOutput(io.MultiWriter(f, os.Stdout))
				} else {
					logrus.SetOutput(f)
				}

				if err := rtLogConf.file.Close(); err != nil {
					panic(err)
				}

				rtLogConf.today = &next
				rtLogConf.file = f
			}

			time.Sleep(500 * time.Millisecond)

			select {
			case <-logExit:
				break loop
			default:
			}
		}
		wg.Done()
	}()
}

// TODO change package
func Time2String(t *time.Time, delim string) string {
	ts := []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
	format := []string{"%d", "%02d", "%02d", "%02d", "%02d", "%02d"}
	tss := make([]string, 6)
	for i, v := range ts {
		tss[i] = fmt.Sprintf(format[i], v)
	}
	return strings.Join(tss, delim)
}

// Stop ...
func Stop() {
	logExit <- nil
	wg.Wait()
	Debug("Logger stop fetching filename with datetime")
}

func getBaseName(fileName string, funcName string) (string, string) {
	return filepath.Base(fileName), filepath.Base(funcName)
}

// Println same as Debug
func Println(args ...interface{}) {
	Debug(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Debug(args...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Debug(args...)
	} else {
		logrus.Debug(args...)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(msg string, args ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Debugf(msg, args...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Debugf(msg, args...)
	} else {
		logrus.Debugf(msg, args...)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Info(args...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Info(args...)
	} else {
		logrus.Info(args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(msg string, args ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Infof(msg, args...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Infof(msg, args...)
	} else {
		logrus.Infof(msg, args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(msg ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Warn(msg...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Warn(msg...)
	} else {
		logrus.Warn(msg...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(msg string, args ...interface{}) {
	if !rtLogConf.showFileInfo {
		logrus.Warnf(msg, args...)
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Warnf(msg, args...)
	} else {
		logrus.Warnf(msg, args...)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(msg ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Error(msg...)
	} else {
		logrus.Error(msg...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(msg string, args ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Errorf(msg, args...)
	} else {
		logrus.Errorf(msg, args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(msg ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Fatal(msg...)
	} else {
		logrus.Fatal(msg...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(msg string, args ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Fatalf(msg, args...)
	} else {
		logrus.Fatalf(msg, args...)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(msg ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Panic(msg...)
	} else {
		logrus.Panic(msg...)
	}
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(msg string, args ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fileName, funcName := getBaseName(file, runtime.FuncForPC(pc).Name())
		logrus.WithField(fileTag, fileName).WithField(lineTag, line).WithField(funcTag, funcName).Panicf(msg, args...)
	} else {
		logrus.Panicf(msg, args...)
	}
}
