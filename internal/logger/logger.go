package logger

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(lokiUrl string) (*zap.Logger, error) {
	// Promtail configuration
	labels := "{source=\"go_app\",job=\"go_app_logger\"}"
	cfg := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.INFO,
	}

	loki, err := promtail.NewClientJson(cfg)
	if err != nil {
		return nil, err
	}

	// Configure Zap to write logs to file and console
	writeSyncer, errorWriter := logWriter()
	if errorWriter != nil {
		return nil, errorWriter
	}

	encoder := getEncoder()

	// Configure core to write to file and console
	logCore := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	log := zap.New(logCore, zap.AddCaller())

	// Configure hooks to send logs to Promtail
	log = log.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		tstamp := time.Now().String()
		loki.Infof(`source = '%s', time = '%s', message = '%s'`, "go app", tstamp, entry.Message)
		time.Sleep(1 * time.Second)
		return nil
	}))

	return log, nil
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		EncodeLevel:  CustomLevelEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func logWriter() (zapcore.WriteSyncer, error) {
	path := "test-log.log"

	logPath, err := createLogFile(path, 0700)
	if err != nil {
		return nil, err
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename: logPath.Name(),
			MaxSize:  500, // MB
			MaxAge:   30,  // days
		}),
		zapcore.AddSync(os.Stdout)), nil
}

func createLogFile(path string, mode fs.FileMode) (*os.File, error) {
	dirName := filepath.Dir(path)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, mode); err != nil {
			return nil, err
		}
	}

	var file *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	err := os.Chmod(path, mode)
	if err != nil {
		return nil, err
	}

	return file, nil
}
