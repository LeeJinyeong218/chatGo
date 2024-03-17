package repository

import (
	"chat/module"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BaseRepository ->
type BaseRepository struct{}

// GetDBCon ->
func (b *BaseRepository) GetDBCon(ctx echo.Context) *gorm.DB {
	dbEngine := ctx.Get(module.CONTEXTDatabase).(*gorm.DB)
	if dbEngine == nil {
		panic("DB is not exist")
	}
	return dbEngine
}
