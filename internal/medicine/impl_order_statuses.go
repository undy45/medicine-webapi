package medicine

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implOrderStatusesApi struct {
}

func NewOrderStatusesApi() OrderStatusesAPI {
	return &implOrderStatusesApi{}
}

func (o implOrderStatusesApi) GetStatus(c *gin.Context) {
	statusId := c.Param("statusId")
	db := HandleConnectionToCollection[Status](c, "db_service_status")
	responseObject, err := db.FindDocument(c, statusId)
	if err != nil {
		HandleRetrievalError(c, err)
		return
	}
	if responseObject != nil {
		c.JSON(http.StatusOK, responseObject)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (o implOrderStatusesApi) GetStatuses(c *gin.Context) {
	db := HandleConnectionToCollection[Status](c, "db_service_status")
	responseObject, err := db.FindAllDocuments(c)
	if err != nil {
		HandleRetrievalError(c, err)
		return
	}
	if responseObject != nil {
		c.JSON(http.StatusOK, responseObject)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
