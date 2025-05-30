package medicine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undy45/medicine-webapi/internal/db_service"
)

func HandleConnectionToCollection[T any](ctx *gin.Context, dbServiceKey string) db_service.DbService[T] {
	value, exists := ctx.Get(dbServiceKey)
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return nil
	}

	db, ok := value.(db_service.DbService[T])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return nil
	}
	return db
}

func HandleRetrievalError(ctx *gin.Context, err error) {
	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load status from database",
				"error":   err.Error(),
			})
		return
	}
}
