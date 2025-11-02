package borrowq

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrowq"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
)

type BorrowQueueHandler interface {
	Enqueue(c echo.Context) error
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
