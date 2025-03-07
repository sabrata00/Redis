package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/kusnadi8605/news/entity"
	"github.com/kusnadi8605/news/usecase"
)

type NewsHandler struct {
	Usecase usecase.NewsUsecase
}

func NewNewsHandler(u usecase.NewsUsecase) *NewsHandler {
	return &NewsHandler{Usecase: u}
}

func (h *NewsHandler) GetAllNews(c echo.Context) error {
	news, err := h.Usecase.GetAllNews(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{Code: 500, Message: "Failed to fetch news", Data: nil})
	}
	return c.JSON(http.StatusOK, entity.Response{Code: 200, Message: "Success", Data: news})
}

func (h *NewsHandler) GetNewsByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{Code: 400, Message: "Invalid ID", Data: nil})
	}

	news, err := h.Usecase.GetNewsByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, entity.Response{Code: 404, Message: "News not found", Data: nil})
	}

	return c.JSON(http.StatusOK, entity.Response{Code: 200, Message: "Success", Data: news})
}

func (h *NewsHandler) CreateNews(c echo.Context) error {
	var news entity.News
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{Code: 400, Message: "Invalid input", Data: nil})
	}

	id, err := h.Usecase.CreateNews(c.Request().Context(), &news)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{Code: 500, Message: "Failed to create news", Data: nil})
	}

	news.ID = int(id)

	return c.JSON(http.StatusCreated, entity.Response{Code: 201, Message: "News created successfully", Data: news})
}

func (h *NewsHandler) UpdateNews(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{Code: 400, Message: "Invalid ID", Data: nil})
	}

	var news entity.News
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{Code: 400, Message: "Invalid input", Data: nil})
	}
	news.ID = id

	err = h.Usecase.UpdateNews(c.Request().Context(), &news)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{Code: 500, Message: "Failed to update news", Data: nil})
	}

	return c.JSON(http.StatusOK, entity.Response{Code: 200, Message: "News updated successfully", Data: news})
}

func (h *NewsHandler) DeleteNews(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{Code: 400, Message: "Invalid ID", Data: nil})
	}

	err = h.Usecase.DeleteNews(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{Code: 500, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, entity.Response{Code: 200, Message: "News deleted successfully", Data: nil})
}
