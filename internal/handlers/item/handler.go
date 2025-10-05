package item

import (
	"net/http"

	itemService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ItemHandler interface {
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
