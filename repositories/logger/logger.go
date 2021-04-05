package logger

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// Logger this is interface for loggers
type Logger interface {
	Log(logType LogType, message string)
}

type BatchLogger interface {
	Log(logType LogType, message []LogMessage)
}

// LogService logging data and using specific logger, logger depend from LogType
type LogService struct {
	tickTime     time.Duration
	loggers      map[LogType]BatchLogger
	MessageQueue map[LogType][]LogMessage
	Mutex        *sync.Mutex
}

func (L *LogService) LaunchLogging() {
	defer func() {
		if err := recover(); err != nil {
			L.Log(Panic, fmt.Sprintf("panic in LaunchLogging func: %v, %v", err, debug.Stack()))
			go L.LaunchLogging()
		}
	}()
	for range time.Tick(L.tickTime) {
		L.Mutex.Lock()
		batches := L.MessageQueue
		L.MessageQueue = make(map[LogType][]LogMessage, 10)
		L.Mutex.Unlock()
		for logType, batch := range batches {
			L.loggers[logType].Log(logType, batch)
		}
	}
}

func (L *LogService) Log(logType LogType, message string) {
	defer func() {
		if msg := recover(); msg != nil {
			L.Log(Panic, fmt.Sprintf("panic in Log func: %v %v", msg, debug.Stack()))
		}
	}()
	L.Mutex.Lock()
	L.MessageQueue[logType] = append(L.MessageQueue[logType], LogMessage{
		Type:    logType,
		Message: message,
	})
	L.Mutex.Unlock()
}

// errors is as
// модели многозадачности
func NewLogService(loggers map[LogType]BatchLogger, period time.Duration) *LogService {
	return &LogService{
		Mutex:        &sync.Mutex{},
		tickTime:     period,
		loggers:      loggers,
		MessageQueue: make(map[LogType][]LogMessage),
	}
}
