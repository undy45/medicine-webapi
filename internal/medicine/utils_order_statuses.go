package medicine

import "github.com/gin-gonic/gin"

type OrderStatusHandler interface {
	GetInitialStatus(c *gin.Context) *Status

	GetStatus(c *gin.Context, statusId string) *Status

	GetStatuses(c *gin.Context) []*Status
}
