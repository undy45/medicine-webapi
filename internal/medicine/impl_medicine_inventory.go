package medicine

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type implMedicineInventoryAPI struct {
}

func NewMedicineInventoryAPI() MedicineInventoryAPI {
	return &implMedicineInventoryAPI{}
}

func (o implMedicineInventoryAPI) DeleteMedicineInventoryEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineInventoryAPI) GetMedicineInventoryEntries(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineInventoryAPI) GetMedicineInventoryEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implMedicineInventoryAPI) UpdateMedicineInventoryEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
