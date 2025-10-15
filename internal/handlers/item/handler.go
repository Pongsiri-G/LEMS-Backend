package item

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	itemService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ItemHandler interface {
	CreateItem(c echo.Context) error
	GetBorrowItem(c echo.Context) error
	GetAll(c echo.Context) error
	GetMyBorrow(c echo.Context) error
	GetChildItemByParentID(c echo.Context) error
	GetFiltered(c echo.Context) error
}

type handler struct {
	service itemService.Service
}

// Consturctor
func NewItemHandler(service itemService.Service) ItemHandler {
	return &handler{service: service}
}

func (h *handler) GetChildItemByParentID(c echo.Context) error {
	itemID := c.Param("item-id")
	response, err := h.service.GetChildItemByParentID(c.Request().Context(), itemID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetBorrowItem(c echo.Context) error {
	itemID := c.Param("item-id")
	item, err := h.service.GetBorrowItem(c.Request().Context(), itemID)

	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
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
		case exceptions.ErrInvalidUUID:
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

func (h *handler) GetAll(c echo.Context) error {
	response, err := h.service.GetAll(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetMyBorrow(c echo.Context) error {
	response, err := h.service.GetMyBorrow(c.Request().Context(), c.Param("user_id"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetFiltered(c echo.Context) error {
	response, err := h.service.GetFiltered(c.Request().Context(), c.Param("strategy"), c.QueryParams()["tags"])

	if err != nil {
		if err == exceptions.ErrNoSuchStrategy {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, response)
}
