/*
 * Medicine Inventory API
 *
 * Medicine inventory management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: your_email@stuba.sk
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package medicine

import (
	"github.com/gin-gonic/gin"
)

type MedicineOrderAPI interface {

	// CreateMedicineOrderEntry Post /api/medicine-order/:ambulanceId/entries
	// Saves new entry into medicine order
	CreateMedicineOrderEntry(c *gin.Context)

	// DeleteMedicineOrderEntry Delete /api/medicine-order/:ambulanceId/entries/:entryId
	// Deletes specific entry
	DeleteMedicineOrderEntry(c *gin.Context)

	// GetMedicineOrderEntries Get /api/medicine-order/:ambulanceId/entries
	// Provides orders of the ambulance
	GetMedicineOrderEntries(c *gin.Context)

	// GetMedicineOrderEntry Get /api/medicine-order/:ambulanceId/entries/:entryId
	// Provides details about ambulance medicine order entry
	GetMedicineOrderEntry(c *gin.Context)

	// UpdateMedicineOrderEntry Put /api/medicine-order/:ambulanceId/entries/:entryId
	// Updates specific entry
	UpdateMedicineOrderEntry(c *gin.Context)
}
