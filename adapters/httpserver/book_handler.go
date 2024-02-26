package httpserver

import (
	"hexagon/adapters/httpserver/model"
	"hexagon/domain/book"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateBook(c echo.Context) error {
	var req model.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	b := book.NewBook(req.ISBN, req.Name)
	if err := s.BookStore.Save(&b); err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) GetBook(c echo.Context) error {
	id := c.Param("id")
	result, err := s.BookStore.FindByISBN(id)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, result)
}

func (s *Server) RegisterBookRoutes(router *echo.Group) {
	router.POST("", s.CreateBook)
	router.GET("/:id", s.GetBook)
}
