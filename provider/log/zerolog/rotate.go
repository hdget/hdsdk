package zerolog

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	defaultBaseDir = "/var/log"
	defaultDirName = "logs"
)

func newRotateLogs(config *RotateLogConf, logFilename string) (*rotatelogs.RotateLogs, error) {
	var rotateLogs *rotatelogs.RotateLogs
	var err error
	switch runtime.GOOS {
	case "linux":
		rotateLogs, err = newLinuxRotateLogs(config, logFilename)
	default:
		rotateLogs, err = newDefaultRotateLogs(config, logFilename)
	}
	return rotateLogs, err
}

func newLinuxRotateLogs(config *RotateLogConf, logFilename string) (*rotatelogs.RotateLogs, error) {
	fileSuffix := path.Ext(logFilename)
	filenameOnly := strings.TrimSuffix(logFilename, fileSuffix)

	// 获取basedir
	basedir := config.BaseDir
	if basedir == "" {
		basedir = defaultBaseDir
	}
	// 创建日志目录
	logDir := path.Join(basedir, filenameOnly)
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}

	logPath := path.Join(logDir, logFilename)
	logFilePrefix := strings.TrimSuffix(logPath, fileSuffix)
	rotateLogFileFormat := logFilePrefix + "%Y%m%d" + fileSuffix
	rotateLogs, _ := rotatelogs.New(
		// 分割后的文件名称
		rotateLogFileFormat,
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(logFilename),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
	)

	return rotateLogs, nil
}

func newDefaultRotateLogs(config *RotateLogConf, logFilename string) (*rotatelogs.RotateLogs, error) {
	dirname := config.Dirname
	if dirname == "" {
		dirname = defaultDirName
	}

	// 创建日志目录
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		return nil, err
	}

	logPath := path.Join(dirname, logFilename)
	fileSuffix := path.Ext(logPath)
	filenameOnly := strings.TrimSuffix(logPath, fileSuffix)
	rotateLogFileFormat := filenameOnly + "%Y%m%d" + fileSuffix
	rotateLogs, _ := rotatelogs.New(
		// 分割后的文件名称
		rotateLogFileFormat,
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
	)

	return rotateLogs, nil
}
