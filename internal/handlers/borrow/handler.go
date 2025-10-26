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
}

type handler struct {
	servicce borrowSvc.Service
}

func NewBorrowHandler(service borrowSvc.Service) BorrowHandler {
	return &handler{servicce: service}
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

	err = h.servicce.Borrow(c.Request().Context(), authUser.ID, req.ItemID)
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

	err = h.servicce.Return(c.Request().Context(), authUser.ID, &req)
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
