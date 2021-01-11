package syslog

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/tchap/zapext/v2/zapsyslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/syslog"
)

var logger *zap.Logger
var ZapcoreLevel zapcore.Level

func InitLogger(logLevel string) error {
	// Flags.
	flagTag := flag.String("syslog_tag", "link-monitor", "syslog tag")

	flag.Parse()

	syslogLevel := syslog.LOG_ERR
	ZapcoreLevel = zapcore.ErrorLevel

	switch logLevel {
	case "debug":
		syslogLevel = syslog.LOG_DEBUG
		ZapcoreLevel = zapcore.DebugLevel
	case "info":
		syslogLevel = syslog.LOG_INFO
		ZapcoreLevel = zapcore.InfoLevel
	case "error":
		syslogLevel = syslog.LOG_ERR
		ZapcoreLevel = zapcore.ErrorLevel
	}

	// Initialize a syslog writer.
	writer, err := syslog.New(syslogLevel|syslog.LOG_LOCAL0, *flagTag)
	if err != nil {
		return errors.Wrap(err, "failed to set up syslog")
	}

	// Initialize Zap.
	encoderConfig := zapcore.EncoderConfig{
		NameKey: "nameKey",
		//TimeKey: "timeKey",
		MessageKey: "something",
		LevelKey: "something",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName: zapcore.FullNameEncoder,
		ConsoleSeparator: " ", // syslog on stock debian doesn't like tabs
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapsyslog.NewCore(ZapcoreLevel, encoder, writer)

	logger = zap.New(core, zap.Development(), zap.AddStacktrace(ZapcoreLevel))
	zap.ReplaceGlobals(logger)

	return errors.Wrap(logger.Sync(), "failed to init logger")
}


func GetLogger() zap.Logger {
	return *logger
}

func Log() *zap.Logger {
	return zap.L()
}
