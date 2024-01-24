package zerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

// 自定义输出格式
func newConsoleLogger() zerolog.ConsoleWriter {
	// 标准输出格式
	w := zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: true,
		// TimeFormat: time.RFC3339,
		TimeFormat: "2006/01/02 15:04:05",
	}

	// format functions
	w.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("msg=\"%s\"", i)
	}

	w.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
		}
		return c
	}

	return w
}
