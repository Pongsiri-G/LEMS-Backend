package tag

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	tagService "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/tag"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TagHandler interface {
	GetNameTagByItemID(c echo.Context) error
	GetTags(c echo.Context) error
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
