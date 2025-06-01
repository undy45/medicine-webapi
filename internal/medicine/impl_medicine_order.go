package medicine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"slices"
)

type implMedicineOrderAPI struct {
}

func NewMedicineOrderAPI() MedicineOrderAPI {
	return &implMedicineOrderAPI{}
}

func (o implMedicineOrderAPI) CreateMedicineOrderEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicineOrderEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.MedicineId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ambulance.MedicineOrders, func(order MedicineOrderEntry) bool {
			return entry.Id == order.Id || entry.MedicineId == order.MedicineId
		})

		statusService := implUtilsOrderStatuses{}
		initialStatus := statusService.GetInitialStatus(c)
		entry.Status = *initialStatus

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.MedicineOrders = append(ambulance.MedicineOrders, entry)
		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.MedicineOrders, func(orderEntry MedicineOrderEntry) bool {
			return entry.Id == orderEntry.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.MedicineOrders[entryIndx], http.StatusOK
	})
}

func (o implMedicineOrderAPI) DeleteMedicineOrderEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineOrders, func(waiting MedicineOrderEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.MedicineOrders = append(ambulance.MedicineOrders[:entryIndx], ambulance.MedicineOrders[entryIndx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

func (o implMedicineOrderAPI) GetMedicineOrderEntries(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.MedicineOrders
		if result == nil {
			result = []MedicineOrderEntry{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

func (o implMedicineOrderAPI) GetMedicineOrderEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.MedicineOrders, func(waiting MedicineOrderEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil ambulance - no need to update it in db
		return nil, ambulance.MedicineOrders[entryIndx], http.StatusOK
	})
}

func (o implMedicineOrderAPI) UpdateMedicineOrderEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry MedicineOrderEntry

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

		entryIndx := slices.IndexFunc(ambulance.MedicineOrders, func(order MedicineOrderEntry) bool {
			return entryId == order.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.Count > 0 {
			ambulance.MedicineOrders[entryIndx].Count = entry.Count
		}

		if entry.MedicineId != "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot update MedicineId in existing order",
			}, http.StatusBadRequest
		}

		if entry.Id != "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot update MedicineId in existing order",
			}, http.StatusBadRequest
		}

		if entry.Name != "" {
			ambulance.MedicineOrders[entryIndx].Name = entry.Name
		}

		if entry.Status.ValidTransitions != nil || entry.Status.Value != "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Can only update status id to change state",
			}, http.StatusBadRequest
		}
		if entry.Status.Id == 0 {
			return ambulance, ambulance.MedicineOrders[entryIndx], http.StatusOK
		}
		currentStatus := ambulance.MedicineOrders[entryIndx].Status
		fmt.Printf("ValidTransitions: %#v (%T), Status.Id: %v (%T)\n", currentStatus.ValidTransitions, currentStatus.ValidTransitions, entry.Status.Id, entry.Status.Id)
		if !slices.Contains(currentStatus.ValidTransitions, entry.Status.Id) {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Changed status is not valid for current order state",
			}, http.StatusBadRequest
		}
		statusService := implUtilsOrderStatuses{}
		changedStatus := statusService.GetStatus(c, int(entry.Status.Id))
		ambulance.MedicineOrders[entryIndx].Status = *changedStatus
		HandleIfDelivered(ambulance, ambulance.MedicineOrders[entryIndx])

		return ambulance, ambulance.MedicineOrders[entryIndx], http.StatusOK
	})
}

func HandleIfDelivered(ambulance *Ambulance, entry MedicineOrderEntry) {
	if entry.Status.Value != "Delivered" {
		return
	}
	foundIndx := slices.IndexFunc(ambulance.MedicineInventory, func(order MedicineInventoryEntry) bool {
		return entry.Id == order.Id || entry.MedicineId == order.MedicineId
	})
	inventoryEntry := ConvertOrderToInventoryEntry(entry)
	if foundIndx >= 0 {
		ambulance.MedicineInventory[foundIndx].Count += entry.Count
	} else {
		ambulance.MedicineInventory = append(ambulance.MedicineInventory, inventoryEntry)
	}

}

func ConvertOrderToInventoryEntry(order MedicineOrderEntry) MedicineInventoryEntry {
	return MedicineInventoryEntry{
		Id:         order.Id,
		MedicineId: order.MedicineId,
		Name:       order.Name,
		Count:      order.Count,
	}
}
