package repository

import (
	"chat/model"
	"github.com/labstack/echo/v4"
)

type (
	ITestRepository interface {
		ListTest(ctx echo.Context) ([]model.TestTable, error)
		GetTest(ctx echo.Context, id int64) (*model.TestTable, error)
		CreateTest(ctx echo.Context, item model.TestTable) (int64, error)
		UpdateTest(ctx echo.Context, id int64, text string) (int64, error)
		DeleteTest(ctx echo.Context, id int64) (int64, error)
	}
	TestRepository struct {
		BaseRepository
	}
)

// NewTestRepository ...
func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func (m *TestRepository) ListTest(ctx echo.Context) ([]model.TestTable, error) {
	db := m.GetDBCon(ctx)

	var items []model.TestTable
	result := db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

func (m *TestRepository) GetTest(ctx echo.Context, id int64) (*model.TestTable, error) {
	db := m.GetDBCon(ctx)

	var item model.TestTable
	result := db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

func (m *TestRepository) CreateTest(ctx echo.Context, item model.TestTable) (int64, error) {
	c := int64(0)

	db := m.GetDBCon(ctx)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := db.Create(&item)
	if result.Error != nil {
		return c, result.Error
	}
	c++

	tx.Commit()

	return c, nil
}

func (m *TestRepository) UpdateTest(ctx echo.Context, id int64, text string) (int64, error) {
	c := int64(0)

	db := m.GetDBCon(ctx)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := db.Where("id = ?", id).Updates(&model.TestTable{Text: text})
	if result.Error != nil {
		return c, result.Error
	}
	c++

	tx.Commit()

	return c, nil
}

func (m *TestRepository) DeleteTest(ctx echo.Context, id int64) (int64, error) {
	c := int64(0)

	db := m.GetDBCon(ctx)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := db.Where("id = ?", id).Delete(&model.TestTable{})
	if result.Error != nil {
		return c, result.Error
	}
	c++

	tx.Commit()

	return c, nil
}
