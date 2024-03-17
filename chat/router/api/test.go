package api

import (
	"chat/dto"
	"chat/module"
	"chat/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TestHandler struct {
	TestService service.ITestService
}

func InitTestHandler(g *echo.Group) {
	testHandler := TestHandler{
		TestService: service.NewTestServiceInstance(),
	}

	// router
	g.GET("", testHandler.listTest)
	g.POST("", testHandler.createTest)
	g.GET("/:id", testHandler.getTest)
	g.PUT("/:id", testHandler.updateTest)
	g.DELETE("/:id", testHandler.deleteTest)

}

func (h *TestHandler) listTest(ctx echo.Context) error {
	items, err := h.TestService.ListTestItem(ctx)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return module.SuccessResponse(ctx, items)
}

func (h *TestHandler) createTest(ctx echo.Context) error {
	var userData dto.TestTableDto
	if err := ctx.Bind(&userData); err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err := h.TestService.CreateTestItem(ctx, userData)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return module.SuccessResponse(ctx, map[string]interface{}{})
}

func (h *TestHandler) getTest(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	item, err := h.TestService.GetTestItem(ctx, id)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return module.SuccessResponse(ctx, item)
}

func (h *TestHandler) updateTest(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	var userData dto.TestTableDto
	if err := ctx.Bind(&userData); err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if id != userData.ID {
		return module.FailedResponse(ctx, http.StatusBadRequest, "id is not equal")
	}

	err = h.TestService.UpdateTestItem(ctx, id, userData)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return module.SuccessResponse(ctx, map[string]interface{}{})

}

func (h *TestHandler) deleteTest(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = h.TestService.DeleteTestItem(ctx, id)
	if err != nil {
		return module.FailedResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return module.SuccessResponse(ctx, map[string]interface{}{})
}
