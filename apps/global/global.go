package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	LOGGER *logrus.Logger
)
