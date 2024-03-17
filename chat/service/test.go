package service

import (
	"chat/dto"
	"chat/model"
	"chat/repository"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"sync"
)

var testInstance *TestService
var testOnce sync.Once

type (
	ITestService interface {
		ListTestItem(ctx echo.Context) ([]dto.TestTableDto, error)
		GetTestItem(ctx echo.Context, id int64) (*dto.TestTableDto, error)
		CreateTestItem(ctx echo.Context, item dto.TestTableDto) error
		UpdateTestItem(ctx echo.Context, id int64, item dto.TestTableDto) error
		DeleteTestItem(ctx echo.Context, id int64) error
	}
	TestService struct {
		testRepository repository.ITestRepository
	}
)

// NewTestServiceInstance ...
func NewTestServiceInstance() *TestService {
	testOnce.Do(func() {
		testInstance = &TestService{
			testRepository: repository.NewTestRepository(),
		}
	})

	return testInstance
}

func (m *TestService) ListTestItem(ctx echo.Context) ([]dto.TestTableDto, error) {
	items, err := m.testRepository.ListTest(ctx)
	if err != nil {
		return nil, err
	}

	var dtoItems []dto.TestTableDto
	for _, item := range items {
		dtoItems = append(dtoItems, dto.TestTableDto{Text: item.Text, ID: item.ID})
	}

	if len(dtoItems) == 0 {
		dtoItems = []dto.TestTableDto{}
	}
	return dtoItems, nil
}

func (m *TestService) GetTestItem(ctx echo.Context, id int64) (*dto.TestTableDto, error) {
	item, err := m.testRepository.GetTest(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v", item)

	return &dto.TestTableDto{
		ID:   item.ID,
		Text: item.Text,
	}, nil
}

func (m *TestService) CreateTestItem(ctx echo.Context, item dto.TestTableDto) error {
	itemModel := model.TestTable{
		ID:   item.ID,
		Text: item.Text,
	}

	c, err := m.testRepository.CreateTest(ctx, itemModel)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("tx error")
	}

	return nil
}

func (m *TestService) UpdateTestItem(ctx echo.Context, id int64, item dto.TestTableDto) error {
	c, err := m.testRepository.UpdateTest(ctx, id, item.Text)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("tx error")
	}

	return nil
}

func (m *TestService) DeleteTestItem(ctx echo.Context, id int64) error {
	c, err := m.testRepository.DeleteTest(ctx, id)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("tx error")
	}

	return nil
}
