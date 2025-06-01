package medicine

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type implOrderStatusesApi struct {
}

func NewOrderStatusesApi() OrderStatusesAPI {
	return &implOrderStatusesApi{}
}

func (o implOrderStatusesApi) GetInitialStatus(c *gin.Context) {
	utilsStatus := implUtilsOrderStatuses{}
	responseObject := utilsStatus.GetInitialStatus(c)
	if responseObject != nil {
		c.JSON(http.StatusOK, responseObject)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (o implOrderStatusesApi) GetStatus(c *gin.Context) {
	statusId := c.Param("statusId")
	utilsStatus := implUtilsOrderStatuses{}
	statusIdInt, _ := strconv.Atoi(statusId)
	responseObject := utilsStatus.GetStatus(c, statusIdInt)
	if responseObject != nil {
		c.JSON(http.StatusOK, responseObject)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (o implOrderStatusesApi) GetStatuses(c *gin.Context) {
	utilsStatus := implUtilsOrderStatuses{}
	responseObject := utilsStatus.GetStatuses(c)
	if responseObject != nil {
		c.JSON(http.StatusOK, responseObject)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
