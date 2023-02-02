package bzerolog

import (
	"github.com/orchestd/log"
	"github.com/orchestd/utils/osutils"
	defaultLogger "github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

type LogSettings struct {
	LogToFile         bool
	FileJsonFormat    bool
	LogToConsole      bool
	ConsoleJsonFormat bool
	CompressLogsFile  bool
	DisableConsoleColor bool
}
func DefaultZeroLogBuilder(conf LogSettings) log.Builder {
	builder := Builder().IncludeCaller()
	if conf.DisableConsoleColor {
		builder = builder.DisableColor()
	}
	var writers []io.Writer
	if conf.LogToFile {
		if conf.FileJsonFormat{
			writers = append(writers,newRollingFile(conf))
		}else{
			writers = append(writers,ConsoleWriter(true,newRollingFile(conf)))
		}
	}
	if conf.LogToConsole {
		if conf.ConsoleJsonFormat {
			writers = append(writers,os.Stderr)
		}else{
			writers = append(writers,ConsoleWriter(conf.DisableConsoleColor,os.Stderr))
		}
	}
	mw := io.MultiWriter(writers...)
	return builder.SetWriter(mw)
}

func newRollingFile(conf LogSettings) io.Writer {
	serviceName := osutils.GetExecutableName()
	if logFolderPath, err := GetLogFolderPath(serviceName); err != nil {
		return nil
	} else if err := os.MkdirAll(logFolderPath, 0744); err != nil {
		defaultLogger.Error().Err(err).Str("path", logFolderPath).Msg("can't create log directory")
		return nil
	} else {

		return &lumberjack.Logger{
			Compress:   conf.CompressLogsFile,
			LocalTime: true,
			Filename:   path.Join(logFolderPath, serviceName+".log"),
			MaxSize:   10,    // megabytes
			MaxAge:     14,     // days
		}
	}
}

func GetLogFolderPath(serviceName string) (string, error) {
	err, path := osutils.GetExecutableDir()
	if err != nil {
		return serviceName, err
	}
	return path + "/Log", nil
}
