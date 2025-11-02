package borrowq

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrowq"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type BorrowQueueHandler interface {
	Enqueue(c echo.Context) error
	MyQueue(c echo.Context) error
	GetFrontQueue(c echo.Context) error
	CancelMyQueue(c echo.Context) error
}

type qorrowQueueHandler struct {
	bqService borrowq.BorrowQueueService
}

func NewBorrowQueueHandler(bqService borrowq.BorrowQueueService) BorrowQueueHandler {
	return &qorrowQueueHandler{
		bqService: bqService,
	}
}

// Enqueue implements BorrowQueueHandler.
func (q *qorrowQueueHandler) Enqueue(c echo.Context) error {
	authUser, err := contextutil.GetUserFromContext(c)

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	var req requests.CreateBorrowQueueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	err = q.bqService.Enqueue(c.Request().Context(), requests.CreateBorrowQueueRequest{
		UserID: authUser.ID,
		ItemID: req.ItemID,
	})

	if err != nil {
		switch err {
		case exceptions.ErrUserAlreadyBorrowq:
			return c.JSON(http.StatusConflict, echo.Map{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "enqueue successfully"})
}

func (q *qorrowQueueHandler) MyQueue(c echo.Context) error {
	itemID := c.Param("itemID")
	if itemID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "required parameter"})
	}

	authUser, err := contextutil.GetUserFromContext(c)

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	res, err := q.bqService.GetGetOneByUserAndItem(c.Request().Context(), itemID, authUser.ID)

	if err != nil {
		switch err {
		case exceptions.ErrBorrowqNotFound:
			return c.JSON(http.StatusNoContent, echo.Map{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, res)
}

// GetFrontQueue implements BorrowQueueHandler.
func (q *qorrowQueueHandler) GetFrontQueue(c echo.Context) error {
    itemID := c.Param("item_id")

    itemUUID, err := uuid.Parse(itemID)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "message": exceptions.ErrInvalidUUID.Error(),
        })
    }

    queue, err := q.bqService.GetFrontQueue(c.Request().Context(), itemUUID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, queue)
}


// CancelMyQueue implements BorrowQueueHandler.
func (q *qorrowQueueHandler) CancelMyQueue(c echo.Context) error {
	queueID := c.Param("queueID")
	if queueID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "required parameter"})
	}

	authUser, err := contextutil.GetUserFromContext(c)

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	err = q.bqService.CancelMyQueue(c.Request().Context(), queueID, authUser.ID)

	if err != nil {
		log.Err(err)
		switch err {
		case exceptions.ErrBorrowqNotFound:
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": exceptions.ErrInternalServer.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK,  map[string]string{ "message": "cancel queue successfully" })	
}