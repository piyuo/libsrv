package log

import "context"

// Logger keep logs
//
type Logger interface {

	// Close Logger
	//
	//	defer logger.Close()
	//
	Close()

	// Write log
	//
	//	logger.write(ctx,"hi","app",DEBUG)
	//
	Write(ctx context.Context, level Level, where, message string)
}
