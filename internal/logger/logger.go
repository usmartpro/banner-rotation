package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/usmartpro/banner-rotation/internal/config"
)

type Logger struct {
	logger *logrus.Logger
}

func New(loggerConfig config.LoggerConf) (*Logger, error) {
	log := logrus.New()

	output, err := openLog(loggerConfig.File)
	if err != nil {
		return nil, fmt.Errorf("error open log file: %w", err)
	}
	log.SetOutput(output)

	level, err := logrus.ParseLevel(loggerConfig.Level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{})

	return &Logger{log}, nil
}

func (l *Logger) Info(message string, params ...interface{}) {
	l.logger.Infof(message, params...)
}

func (l *Logger) Error(message string, params ...interface{}) {
	l.logger.Errorf(message, params...)
}

func (l *Logger) LogRequest(r *http.Request, code, length int) {
	l.logger.Infof(
		"%s %s %s %s %d %d %q",
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

func openLog(file string) (io.Writer, error) {
	switch file {
	case "stderr":
		fmt.Println("stderr")
		return os.Stderr, nil
	case "stdout":
		fmt.Println("stdout")
		return os.Stdout, nil
	default:
		file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
