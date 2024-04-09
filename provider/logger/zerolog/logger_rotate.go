package zerolog

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func newRotateLogger(conf *zerologProviderConfig) (io.Writer, error) {
	if conf.Filename == "" {
		return nil, errdef.ErrInvalidConfig
	}

	// 获取logDir
	dir := conf.Dir
	if dir == "" {
		switch runtime.GOOS {
		case "linux":
			dir = linuxDefaultDir
		default:
			dir = nonLinuxDefaultDir
		}
	}

	// 创建日志目录
	fileSuffix := path.Ext(conf.Filename)
	filename := strings.TrimSuffix(conf.Filename, fileSuffix)
	rotateDir := path.Join(dir, filename)
	err := os.MkdirAll(rotateDir, 0744)
	if err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(rotateDir, conf.Filename),
		MaxSize:    conf.Rotate.MaxSize,   // The maximum size in megabytes of the log file before it gets rotated, It defaults to 100 megabytes.
		MaxAge:     conf.Rotate.MaxAge,    // In days before deleting the file
		MaxBackups: conf.Rotate.MaxBackup, // The maximum number of old log files to retain.
		Compress:   conf.Rotate.Compress,  // Compress the rotated log files, false by default.
	}, nil
}
