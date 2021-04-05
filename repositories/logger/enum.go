package logger

// LogType represent type of logs
type LogType = string

// represent values of LogType
const (
	Error    LogType = "error"
	Panic    LogType = "panic"
	Incoming LogType = "incoming"
)

// LogMessage format message for logging
type LogMessage struct {
	Type    LogType `json:"type"`
	Message string  `json:"message"`
}
