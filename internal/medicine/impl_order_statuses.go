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

func (o implOrderStatusesApi) GetStatuses(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
