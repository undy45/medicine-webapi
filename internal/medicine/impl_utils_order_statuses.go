package medicine

import (
	"github.com/gin-gonic/gin"
)

type implUtilsOrderStatuses struct {
}

func UtilsOrderStatuses() OrderStatusHandler {
	return &implUtilsOrderStatuses{}
}

func (o implUtilsOrderStatuses) GetInitialStatus(c *gin.Context) *Status {
	statusId := 1
	db := HandleConnectionToCollection[Status](c, "db_service_status")
	responseObject, err := db.FindDocument(c, statusId)
	if err != nil {
		HandleRetrievalError(c, err)
		return nil
	}
	return responseObject
}

func (o implUtilsOrderStatuses) GetStatus(c *gin.Context, statusId int) *Status {
	db := HandleConnectionToCollection[Status](c, "db_service_status")
	responseObject, err := db.FindDocument(c, statusId)
	if err != nil {
		HandleRetrievalError(c, err)
		return nil
	}
	return responseObject
}

func (o implUtilsOrderStatuses) GetStatuses(c *gin.Context) []*Status {
	db := HandleConnectionToCollection[Status](c, "db_service_status")
	responseObject, err := db.FindAllDocuments(c)
	if err != nil {
		HandleRetrievalError(c, err)
		return nil
	}
	return responseObject
}
