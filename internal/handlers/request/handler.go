package request

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	requestSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/request"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestHandler interface {
	GetRequests(c echo.Context) error
	GetMyRequests(c echo.Context) error

	CreateRequest(c echo.Context) error
	EditRequest(c echo.Context) error
}

type handler struct {
	service requestSvc.Service
}

func NewRequestHandler(service requestSvc.Service) RequestHandler {
	return &handler{service: service}
}

// CreateRequest implements RequestHandler.
func (h *handler) CreateRequest(c echo.Context) error {
	var req requests.CreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "bad request",
		})
	}

	auth, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	userID, err := uuid.Parse(auth.ID)
	if err != nil {
		return c.JSON(401, nil)
	}

	err = h.service.CreateRequest(c.Request().Context(), &userID, req)
	if err != nil {
		switch err {
		case exceptions.ErrRequestNotExpectItemID:
		case exceptions.ErrRequestNotExpectItem:
		case exceptions.ErrRequestInvalidRequestType:
			return c.JSON(400, nil)
		case exceptions.ErrInvalidUUID:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrInvalidUUID.Error(),
			})
		case exceptions.ErrItemNotFound:
			return c.JSON(404, nil)
		case exceptions.ErrRequestItemInvalid:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrRequestItemInvalid.Error(),
			})
		case exceptions.ErrRequestItemIDInvalid:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrRequestItemIDInvalid.Error(),
			})
		case exceptions.ErrUserNotFound:
			return c.JSON(404, nil)
		case exceptions.ErrUserIDIsNil:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrUserIDIsNil.Error(),
			})
		default:
			return c.JSON(500, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(201, nil)
}

// EditRequest implements RequestHandler.
func (h *handler) EditRequest(c echo.Context) error {
	var req requests.EditRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "bad request",
		})
	}

	err = h.service.EditRequest(c.Request().Context(), req)
	if err != nil {
		switch err {
		case exceptions.ErrRequestInvalidRequestType:
		case exceptions.ErrRequestNotFound:
		case exceptions.ErrItemNotFound:
		case exceptions.ErrUserNotFound:
			return c.JSON(404, nil)
		case exceptions.ErrInvalidUUID:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrInvalidUUID.Error(),
			})
		case exceptions.ErrRequestNotFound:
			return c.JSON(404, nil)
		default:
			return c.JSON(500, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}
	return c.JSON(200, nil)
}

// GetMyRequests implements RequestHandler.
func (h *handler) GetMyRequests(c echo.Context) error {
	auth, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	userID, err := uuid.Parse(auth.ID)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": exceptions.ErrInvalidUUID.Error(),
		})
	}

	requestsData, err := h.service.GetRequests(c.Request().Context(), &userID)
	if err != nil {
		switch err {
		case exceptions.ErrItemNotFound:
		case exceptions.ErrUserNotFound:
			return c.JSON(404, nil)
		case exceptions.ErrInvalidUUID:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrInvalidUUID.Error(),
			})
		default:
			return c.JSON(500, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}

	return c.JSON(200, requestsData)
}

// GetRequests implements RequestHandler.
func (h *handler) GetRequests(c echo.Context) error {
	requestsData, err := h.service.GetRequests(c.Request().Context(), nil)
	if err != nil {
		switch err {
		case exceptions.ErrItemNotFound:
		case exceptions.ErrUserNotFound:
			return c.JSON(404, nil)
		case exceptions.ErrInvalidUUID:
			return c.JSON(400, echo.Map{
				"message": exceptions.ErrInvalidUUID.Error(),
			})
		default:
			return c.JSON(500, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}

	return c.JSON(200, requestsData)
}
