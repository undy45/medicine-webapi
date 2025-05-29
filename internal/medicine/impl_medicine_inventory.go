package medicine

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

type implMedicineInventoryAPI struct {
}

func NewMedicineInventoryAPI() MedicineInventoryAPI {
	return &implMedicineInventoryAPI{}
}

func (o implMedicineInventoryAPI) DeleteMedicineInventoryEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineInventory, func(waiting MedicineInventoryEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.MedicineInventory = append(ambulance.MedicineInventory[:entryIndx], ambulance.MedicineInventory[entryIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

func (o implMedicineInventoryAPI) GetMedicineInventoryEntries(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.MedicineInventory
		if result == nil {
			result = []MedicineInventoryEntry{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

func (o implMedicineInventoryAPI) GetMedicineInventoryEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineInventory, func(waiting MedicineInventoryEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.MedicineInventory[entryIndx], http.StatusOK
	})
}

func (o implMedicineInventoryAPI) UpdateMedicineInventoryEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicineInventoryEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineInventory, func(inventory MedicineInventoryEntry) bool {
			return entryId == inventory.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.Count > 0 {
			ambulance.MedicineInventory[entryIndx].Count = entry.Count
		} else if entry.Count == 0 {
			ambulance.MedicineInventory = append(ambulance.MedicineInventory[:entryIndx], ambulance.MedicineInventory[entryIndx+1:]...)
			return ambulance, nil, http.StatusOK
		}

		if entry.MedicineId != "" {
			ambulance.MedicineInventory[entryIndx].MedicineId = entry.MedicineId
		}

		if entry.Id != "" {
			ambulance.MedicineInventory[entryIndx].Id = entry.Id
		}

		if entry.Name != "" {
			ambulance.MedicineInventory[entryIndx].Name = entry.Name
		}

		return ambulance, ambulance.MedicineInventory[entryIndx], http.StatusOK
	})
}
