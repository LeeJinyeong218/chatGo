package module

import (
	"chat/dto"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MiddlewareEnvironments -> 앱 실행시 context 등록
func MiddlewareEnvironments(config *dto.Config, dbEngine *gorm.DB, loggerEntry *logrus.Entry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(CONTEXTConfig, config)
			c.Set(CONTEXTDatabase, dbEngine)
			c.Set(CONTEXTLogger, loggerEntry)
			return next(c)
		}
	}
}
