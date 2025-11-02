package borrow

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	borrowSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type BorrowHandler interface {
	Return(c echo.Context) error
	Borrow(c echo.Context) error
	GetMyBorrowLog(c echo.Context) error
	GetBorrowLog(c echo.Context) error
	GetBorrowID(c echo.Context) error
}

type handler struct {
	service borrowSvc.Service
}

func NewBorrowHandler(service borrowSvc.Service) BorrowHandler {
	return &handler{service: service}
}

func (h *handler) Borrow(c echo.Context) error {
	var req requests.BorrowRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	err = h.service.Borrow(c.Request().Context(), authUser.ID, req.ItemID)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrItemQuantityInSufficient:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": exceptions.ErrItemQuantityInSufficient.Error(),
			})
		case exceptions.ErrFailedToUpdateQuantity:
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrFailedToUpdateQuantity.Error(),
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		case exceptions.ErrNotYourTurnInQueue:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": exceptions.ErrNotYourTurnInQueue.Error(),
			})
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "borrow item successfully",
	})
}

// Return implements Handler.
func (h *handler) Return(c echo.Context) error {
	var req requests.ReturnRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	err = h.service.Return(c.Request().Context(), authUser.ID, &req)
	if err != nil {
		switch err {
		case exceptions.ErrInvalidUUID:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid uuid format",
			})
		case exceptions.ErrBorrowLogNotFound:
			return c.JSON(http.StatusNotFound, nil)
		case exceptions.ErrItemNotFound:
			return c.JSON(http.StatusNotFound, nil)
		case exceptions.ErrCannotReturnChildItemDirectly:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": exceptions.ErrCannotReturnChildItemDirectly.Error(),
			})
		default:
			log.Error().Err(err).Msg(exceptions.ErrInternalServer.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "return item successfully",
	})
}

// GetMyBorrowLog implements BorrowHandler.
func (h *handler) GetMyBorrowLog(c echo.Context) error {
	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	borrowLogs, err := h.service.GetUsersBorrowedItems(c.Request().Context(), authUser.ID)
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
	return c.JSON(http.StatusOK, borrowLogs)
}

// GetBorrowLog implements BorrowHandler.
func (h *handler) GetBorrowLog(c echo.Context) error {
	borrowLogs, err := h.service.GetAllBorrowedItems(c.Request().Context())
	if err != nil {
		switch err {
		default:
			log.Error().Err(err).Msg("error while getting all borrowed items")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, borrowLogs)

}

func (h *handler) GetBorrowID(c echo.Context) error {
	userID, err := contextutil.GetUserFromContext(c)

	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	itemID := c.Param("item-id")

	borrowID, err := h.service.GetBorrowID(c.Request().Context(), userID.ID, itemID)

	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, borrowID)

}
