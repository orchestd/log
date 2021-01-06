package log

import (
	"context"

)

const (
	compensateMortarLoggerWrapper = 1
)

type loggerWrapper struct {
	contextExtractors []ContextExtractor
	logger            Logger
}

// CreateMortarLogger creates a new mortar logger which is a wrapper to support
// 	- ContextExtractors
//
// **Important**
//	This constructor will call builder.IncrementSkipFrames to peel additional layer of itself.
func CreateMortarLogger(builder Builder, contextExtractors ...ContextExtractor) Logger {
	logger := builder.IncrementSkipFrames(compensateMortarLoggerWrapper).Build() // add 1
	return &loggerWrapper{
		contextExtractors: contextExtractors,
		logger:            logger,
	}
}

func (l *loggerWrapper) Trace(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Trace(ctx, format, args...)
}

func (l *loggerWrapper) Debug(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Debug(ctx, format, args...)
}

func (l *loggerWrapper) Info(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Info(ctx, format, args...)
}

func (l *loggerWrapper) Warn(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Warn(ctx, format, args...)
}

func (l *loggerWrapper) Error(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Error(ctx, format, args...)
}

func (l *loggerWrapper) Custom(ctx context.Context, level Level, skipAdditionalFrames int, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Custom(ctx, level, skipAdditionalFrames, format, args...)
}

func (l *loggerWrapper) WithError(err error) Fields {
	return newEntry(l.contextExtractors, l.logger, true).WithError(err)
}

func (l *loggerWrapper) WithField(name string, value interface{}) Fields {
	return newEntry(l.contextExtractors, l.logger, true).WithField(name, value)
}

func (l *loggerWrapper) WithFields(fields map[string]interface{}) Fields {
	return newEntry(l.contextExtractors, l.logger, true).WithFields(fields)
}

func (l *loggerWrapper) Configuration() LoggerConfiguration {
	return l.logger.Configuration()
}
