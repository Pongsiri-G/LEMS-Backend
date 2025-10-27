package tag

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	tagService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/tag"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TagHandler interface {
	GetNameTagByItemID(c echo.Context) error
	GetTags(c echo.Context) error
	CreateTag(c echo.Context) error
	UnAssignTagFromItem(c echo.Context) error
}

type handler struct {
	service tagService.Service
}

// Consturctor
func NewTagHandler(service tagService.Service) TagHandler {
	return &handler{service: service}
}

func (h *handler) GetNameTagByItemID(c echo.Context) error {
	itemID := c.Param("itemID")
	tags, err := h.service.GetTagsNameByItemID(c.Request().Context(), itemID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid param",
		})
	}
	return c.JSON(http.StatusOK, tags)
}

// GetTags implements TagHandler.
func (h *handler) GetTags(c echo.Context) error {
	tags, err := h.service.GetAllTags(c.Request().Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get tags")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": exceptions.ErrInternalServer.Error(),
		})
	}
	return c.JSON(http.StatusOK, tags)
}

// CreateTag implements TagHandler.
func (h *handler) CreateTag(c echo.Context) error {
	var req requests.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind create tag request")
		return c.JSON(http.StatusBadRequest, nil)
	}

	response, err := h.service.CreateTag(c.Request().Context(), &req)
	if err != nil {
		switch err {
		default:
			log.Error().Err(err).Msg("failed to create tag")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusCreated, response)
}

// UnAssignTagFromItem implements TagHandler.
func (h *handler) UnAssignTagFromItem(c echo.Context) error {
	itemID := c.Param("item_id")
	tagID := c.Param("tag_id")

	err := h.service.UnAssignTagFromItem(c.Request().Context(), itemID, tagID)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrInvalidUUID.Error(),
			})
		case exceptions.ErrTagNotAssigned:
			return c.JSON(404, nil)
		default:
			log.Error().Err(err).Msg("failed to unassign tag from item")
			return c.JSON(500, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(200, nil)
}
