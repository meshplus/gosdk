package common

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

const defaultLoggerFormat = "%{color}[%{module}][%{level:.5s}] %{time:2006-01-02 15:04:05} %{shortfile} %{message} %{color:reset}"

var (
	format    = logging.MustStringFormatter(defaultLoggerFormat)
	loggers   = make(map[string]*logging.Logger)
	logger    = logging.MustGetLogger("common")
	backend   logging.LeveledBackend
	conf      *viper.Viper
	once      sync.Once
	lock      sync.RWMutex
	autoReset = true
)

func newConsoleBackend(vip *viper.Viper) logging.LeveledBackend {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)

	// set level
	logLevel := vip.GetString(LogOutputLevel)

	logger.Debugf("[CONFIG]: %s = %v", LogOutputLevel, logLevel)

	level, _ := logging.LogLevel(logLevel)

	backendLeveled.SetLevel(level, "")

	return backendLeveled
}

func newFileBackend(vip *viper.Viper) logging.LeveledBackend {
	dir := vip.GetString(LogDir)
	logger.Debugf("[CONFIG]: %s = %v", LogDir, dir)

	fileName := path.Join(dir, "gosdk"+time.Now().Format("-2006-01-02-15:04:05 PM")+".log")
	_ = os.MkdirAll(dir, 0777)
	file, err := os.Create(fileName)
	if err != nil {
		logger.Errorf("create %s failed", fileName)
	}
	fileBackend := logging.NewLogBackend(file, "", 0)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, format)
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)

	// set level
	logLevel := vip.GetString(LogOutputLevel)

	logger.Debugf("[CONFIG]: %s = %v", LogOutputLevel, logLevel)

	level, _ := logging.LogLevel(logLevel)

	fileBackendLeveled.SetLevel(level, "")

	return fileBackendLeveled
}

func Reset() {
	autoReset = true
}

func updateBackend() {
	t := time.NewTicker(24 * time.Hour)
	if autoReset {
		for range t.C {
			SetBackend()
		}
	}
}

func InitLog(vip *viper.Viper) {
	once.Do(func() {
		conf = vip
		SetBackend()
		go updateBackend()
	})
}

func SetBackend() {
	consoleBackendLeveled := newConsoleBackend(conf)
	fileBackendLeveled := newFileBackend(conf)

	backend = logging.MultiLogger(consoleBackendLeveled, fileBackendLeveled)
	lock.Lock()
	for _, logger := range loggers {
		logger.SetBackend(backend)
	}
	lock.Unlock()
}

func SetCustomBackend(customBackend logging.LeveledBackend) {
	lock.Lock()
	defer lock.Unlock()
	autoReset = false
	for _, logger := range loggers {
		logger.SetBackend(customBackend)
	}
}

func GetLogger(module string) *logging.Logger {
	var logger *logging.Logger

	lock.Lock()
	if loggers[module] != nil {
		return loggers[module]
	} else {
		logger = logging.MustGetLogger(module)
		if backend != nil {
			logger.SetBackend(backend)
		}
		loggers[module] = logger
	}
	lock.Unlock()

	return logger
}
