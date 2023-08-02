package initialize

import (
	"easy-tiktok/apps/global"
	"github.com/sirupsen/logrus"
	"os"
)

// init //
// 初始化logrus的日志类
// Author lql
func init() {
	// 初始化新Logger示例
	logger := logrus.New()
	// 设置输出
	logger.Out = os.Stdout
	// 设置格式
	logger.Formatter = &logrus.JSONFormatter{}
	// 设置输出级别
	logger.SetLevel(logrus.InfoLevel)
	global.LOGGER = logger
}
