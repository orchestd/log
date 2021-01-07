package log

import (
	"context"
)

const (
	noAdditionalSkipFrames = 0
)

type contextAwareLogEntry struct {
	contextExtractors []ContextExtractor
	innerLogger       Fields
	fields            map[string]interface{}
	err               error
	withFields        bool
}

func newEntry(contextExtractors []ContextExtractor, logger Fields, withFields bool) Fields {
	return &contextAwareLogEntry{
		contextExtractors: contextExtractors,
		innerLogger:       logger,
		fields:            make(map[string]interface{}),
		err:               nil,
		withFields:        withFields,
	}
}

func (c *contextAwareLogEntry) Trace(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, TraceLevel, noAdditionalSkipFrames, format, args...)
}

func (c *contextAwareLogEntry) Debug(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, DebugLevel, noAdditionalSkipFrames, format, args...)
}

func (c *contextAwareLogEntry) Info(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, InfoLevel, noAdditionalSkipFrames, format, args...)
}

func (c *contextAwareLogEntry) Warn(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, WarnLevel, noAdditionalSkipFrames, format, args...)
}

func (c *contextAwareLogEntry) Error(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, ErrorLevel, noAdditionalSkipFrames, format, args...)
}

func (c *contextAwareLogEntry) Custom(ctx context.Context, level Level, skipAdditionalFrames int, format string, args ...interface{}) {
	c.log(ctx, level, skipAdditionalFrames, format, args...)
}

func (c *contextAwareLogEntry) WithError(err error) Fields {
	c.err = err
	return c
}

func (c *contextAwareLogEntry) WithField(name string, value interface{}) Fields {
	c.fields[name] = value
	return c
}

func (c *contextAwareLogEntry) WithFields(fields map[string]interface{}) Fields {
	for name, value := range fields {
		c.fields[name] = value
	}
	return c
}

func (c *contextAwareLogEntry) log(ctx context.Context, level Level, skipAdditionalFrames int, format string, args ...interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}
	logger := c.enrich(ctx)
	for k, v := range c.fields {
		logger = logger.WithField(k, v)
	}
	if c.err != nil {
		logger = logger.WithError(c.err)
	}
	if !c.withFields { // if no fields, we have one less layer to peel
		skipAdditionalFrames++
	}
	logger.Custom(ctx, level, skipAdditionalFrames, format, args...)
}

func (c *contextAwareLogEntry) enrich(ctx context.Context) (logger Fields) {
	defer func() {
		if r := recover(); r != nil {
			c.innerLogger.WithField("__panic__", r).Error(ctx, "one of the context extractors panicked")
			logger = c.innerLogger
		}
	}()
	logger = c.innerLogger
	for _, extractor := range c.contextExtractors {
		for k, v := range extractor(ctx) {
			logger = logger.WithField(k, v)
		}
	}
	return
}
