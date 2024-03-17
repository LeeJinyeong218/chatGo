package router

import (
	"chat/dto"
	"chat/module"
	"chat/router/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Router -> echo middleware 및 라우팅 설정
func Router(config *dto.Config, dbEngine *gorm.DB, loggerEntry *logrus.Entry) {
	e := echo.New()

	e.Use(module.MiddlewareEnvironments(config, dbEngine, loggerEntry))

	api.InitTestHandler(e.Group("/test"))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Server.BaseUrl, config.Server.Port)))
}
