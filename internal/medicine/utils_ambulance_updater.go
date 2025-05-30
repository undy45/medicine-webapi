package medicine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undy45/medicine-webapi/internal/db_service"
)

type ambulanceUpdater = func(
	ctx *gin.Context,
	ambulance *Ambulance,
) (updatedAmbulance *Ambulance, responseContent interface{}, status int)

func updateAmbulanceFunc(ctx *gin.Context, updater ambulanceUpdater) {
	ambulanceId := ctx.Param("ambulanceId")
	db := HandleConnectionToCollection[Ambulance](ctx, "db_service_ambulance")
	ambulance, err := db.FindDocument(ctx, ambulanceId)
	if err != nil {
		HandleRetrievalError(ctx, err)
		return
	}

	updatedAmbulance, responseObject, status := updater(ctx, ambulance)

	if updatedAmbulance != nil {
		err = db.UpdateDocument(ctx, ambulanceId, updatedAmbulance)
	} else {
		err = nil // redundant but for clarity
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update ambulance in database",
				"error":   err.Error(),
			})
	}

}
