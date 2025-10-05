package item

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	itemService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ItemHandler interface {
	CreateItem(c echo.Context) error
	GetBorrowItem(c echo.Context) error
}

type handler struct {
	service itemService.Service
}

// Consturctor
func NewItemHandler(service itemService.Service) ItemHandler {
	return &handler{service: service}
}

func (h *handler) GetBorrowItem(c echo.Context) error {
	itemID := c.Param("itemID")
	item, err := h.service.GetBorrowItem(c.Request().Context(), itemID)

	if err != nil {
		switch err {
		case itemService.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		default:
			log.Error().Err(err).Msg("internal server error")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "internal server error",
			})
		}
	}
	return c.JSON(http.StatusOK, item)

}

// CreateItem implements ItemHandler.
func (h *handler) CreateItem(c echo.Context) error {
	var req requests.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}
	if err := h.service.CreateItem(c.Request().Context(), &req); err != nil {
		switch err {
		case itemService.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		default:
			log.Error().Err(err).Msg("internal server error")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "internal server error",
			})
		}
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "item created successfully",
	})
}
