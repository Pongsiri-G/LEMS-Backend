package item

import (
	"fmt"
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	itemService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ItemHandler interface {
	CreateItem(c echo.Context) error
	GetBorrowItem(c echo.Context) error
	GetAll(c echo.Context) error
	GetMyBorrow(c echo.Context) error
	GetChildItemByParentID(c echo.Context) error
	SearchItems(c echo.Context) error
	DeleteItem(c echo.Context) error

	AssignItemSet(c echo.Context) error
	RemoveItemSet(c echo.Context) error
	EditItem(c echo.Context) error
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
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
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
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
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
			"message": exceptions.ErrInternalServer.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetMyBorrow(c echo.Context) error {
	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}
	response, err := h.service.GetMyBorrow(c.Request().Context(), authUser.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": exceptions.ErrInternalServer.Error(),
		})
	}

	fmt.Println("HAAHAH", response)
	return c.JSON(http.StatusOK, response)
}

func (h *handler) SearchItems(c echo.Context) error {
	name := c.QueryParam("name")
	tags := c.QueryParams()["tags"]
	status := c.QueryParam("status")
	user := c.QueryParam("user")
	var userId = ""

	if user != "" {
		userDetail, err := contextutil.GetUserFromContext(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Internal Server Error",
			})
		}
		userId = userDetail.ID
	}

	strategies := item.SearchStrategyMap{
		Name:   name,
		Tags:   tags,
		Status: status,
		User:   userId,
	}

	response, err := h.service.SearchItems(c.Request().Context(), strategies)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(http.StatusOK, response)
}

// DeleteItem implements ItemHandler.
func (h *handler) DeleteItem(c echo.Context) error {
	itemID := c.Param("item-id")
	err := h.service.DeleteItem(c.Request().Context(), itemID)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, nil)
}

// AssignItemSet implements ItemHandler.
func (h *handler) AssignItemSet(c echo.Context) error {
	parentID := c.Param("parent-item-id")
	childID := c.Param("child-item-id")

	err := h.service.AssignChild(c.Request().Context(), parentID, childID)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		case exceptions.ErrItemSetAlreadyExists:
			return c.JSON(http.StatusConflict, echo.Map{
				"message": "item set already exists",
			})
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusCreated, nil)
}

// EditItem implements ItemHandler.
func (h *handler) EditItem(c echo.Context) error {
	var req requests.EditItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}
	err := h.service.UpdateItem(c.Request().Context(), &req)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "item edited successfully",
	})
}

// RemoveItemSet implements ItemHandler.
func (h *handler) RemoveItemSet(c echo.Context) error {
	parentID := c.Param("parent-item-id")
	childID := c.Param("child-item-id")

	err := h.service.RemoveChild(c.Request().Context(), parentID, childID)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, nil)
}
