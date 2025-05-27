package medicine

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type implMedicineOrderAPI struct {
}

func NewMedicineOrderAPI() MedicineOrderAPI {
	return &implMedicineOrderAPI{}
}

func (o implMedicineOrderAPI) CreateMedicineOrderEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineOrderAPI) DeleteMedicineOrderEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineOrderAPI) GetMedicineOrderEntries(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineOrderAPI) GetMedicineOrderEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineOrderAPI) UpdateMedicineOrderEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
