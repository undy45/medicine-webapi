package medicine

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/undy45/medicine-webapi/internal/db_service"
)

type MedicineOrderSuite struct {
	suite.Suite
	dbAmbulanceServiceMock *DbServiceMock[Ambulance]
	dbStatusServiceMock    *DbServiceMock[Status]
}

func TestMedicineOrderSuite(t *testing.T) {
	suite.Run(t, new(MedicineOrderSuite))
}

func (suite *MedicineOrderSuite) SetupTest() {
	suite.dbAmbulanceServiceMock = &DbServiceMock[Ambulance]{}
	suite.dbStatusServiceMock = &DbServiceMock[Status]{}

	// Compile time Assert that the mock is of type db_service.DbService[Ambulance]
	var _ db_service.DbService[Ambulance] = suite.dbAmbulanceServiceMock
	var _ db_service.DbService[Status] = suite.dbStatusServiceMock

	suite.dbAmbulanceServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id: "test-ambulance",
				MedicineOrders: []MedicineOrderEntry{
					{
						Id:         "test-entry",
						Name:       "test-name",
						MedicineId: "test-medicine-id",
						Count:      15,
						Status: Status{
							Id:               1,
							Value:            "To_ship",
							ValidTransitions: []int32{2, 4},
						},
					},
				},
			},
			nil,
		)

	suite.dbStatusServiceMock.
		On("FindDocument", mock.Anything, "1").
		Return(
			&Status{
				Id:               1,
				Value:            "To_ship",
				ValidTransitions: []int32{2, 4},
			},
			nil,
		).
		On("FindDocument", mock.Anything, "2").
		Return(
			&Status{
				Id:               2,
				Value:            "Shipped",
				ValidTransitions: []int32{3, 4},
			},
			nil,
		).
		On("FindDocument", mock.Anything, "3").
		Return(
			&Status{
				Id:               3,
				Value:            "Delivered",
				ValidTransitions: []int32{},
			},
			nil,
		).
		On("FindDocument", mock.Anything, "4").
		Return(
			&Status{
				Id:               4,
				Value:            "Canceled",
				ValidTransitions: []int32{},
			},
			nil,
		)
	suite.dbStatusServiceMock.
		On("FindAllDocuments", mock.Anything, mock.Anything).
		Return(
			[]*Status{
				{
					Id:               1,
					Value:            "To_ship",
					ValidTransitions: []int32{2, 4},
				},
				{
					Id:               2,
					Value:            "Shipped",
					ValidTransitions: []int32{3, 4},
				},
				{
					Id:               3,
					Value:            "Delivered",
					ValidTransitions: []int32{},
				},
				{
					Id:               4,
					Value:            "Canceled",
					ValidTransitions: []int32{},
				},
			},
			nil,
		)

	suite.dbAmbulanceServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func (suite *MedicineOrderSuite) Test_DeleteOrder_DbService() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("DELETE", "/medicine-order/test-ambulance/entries/test-entry", nil)

	sut := implMedicineOrderAPI{}

	// ACT
	sut.DeleteMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusNoContent, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			return len(arg.MedicineOrders) == 0
		}),
	)
}

func (suite *MedicineOrderSuite) Test_GetOrder_DbService() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("GET", "/medicine-order/test-ambulance/entries/test-entry", nil)

	sut := implMedicineOrderAPI{}

	// ACT
	sut.GetMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	var respObj map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	suite.Equal("test-entry", respObj["id"])
	suite.Equal("test-name", respObj["name"])
	suite.Equal("test-medicine-id", respObj["medicineId"])
	suite.Equal(float64(15), respObj["count"])
	var respStatus = respObj["status"].(map[string]interface{})
	suite.Equal(float64(1), respStatus["id"])
	suite.Equal("To_ship", respStatus["value"])
	var respValidTransitionsIface = respStatus["validTransitions"].([]interface{})
	respValidTransitions := make([]int32, len(respValidTransitionsIface))
	for i, v := range respValidTransitionsIface {
		respValidTransitions[i] = int32(v.(float64))
	}
	suite.ElementsMatch([]int32{2, 4}, respValidTransitions)
	suite.dbAmbulanceServiceMock.AssertNotCalled(suite.T(), "UpdateDocument", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *MedicineOrderSuite) Test_GetOrder_DbServiceGetAllEntries() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("GET", "/medicine-order/test-ambulance/entries/", nil)

	sut := implMedicineOrderAPI{}

	// ACT
	sut.GetMedicineOrderEntries(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	log.Println(recorder.Body.String())
	var respObj []map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	suite.Equal(len(respObj), 1)
	var respOrder = respObj[0]
	suite.Equal("test-entry", respOrder["id"])
	suite.Equal("test-name", respOrder["name"])
	suite.Equal("test-medicine-id", respOrder["medicineId"])
	suite.Equal(float64(15), respOrder["count"])
	var respStatus = respOrder["status"].(map[string]interface{})
	suite.Equal(float64(1), respStatus["id"])
	suite.Equal("To_ship", respStatus["value"])
	var respValidTransitionsIface = respStatus["validTransitions"].([]interface{})
	respValidTransitions := make([]int32, len(respValidTransitionsIface))
	for i, v := range respValidTransitionsIface {
		respValidTransitions[i] = int32(v.(float64))
	}
	suite.ElementsMatch([]int32{2, 4}, respValidTransitions)
	suite.dbAmbulanceServiceMock.AssertNotCalled(suite.T(), "UpdateDocument", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *MedicineOrderSuite) Test_CreateOrder_DbService() {
	// ARRANGE
	json := `{
        "id": "input-entry-id",
        "name": "input-test-name",
        "medicineId": "input-test-medicine-id",
		"count": 20
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("POST", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.CreateMedicineOrderEntry(ctx)

	// ASSERT
	var expectedObj = &MedicineOrderEntry{
		Id:         "input-entry-id",
		Name:       "input-test-name",
		MedicineId: "input-test-medicine-id",
		Count:      20,
		Status: Status{
			Id:               1,
			Value:            "To_ship",
			ValidTransitions: []int32{2, 4},
		},
	}
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			for _, entry := range arg.MedicineOrders {
				if reflect.DeepEqual(entry, *expectedObj) {
					return true
				}
			}
			return false
		}),
	)
}

func (suite *MedicineOrderSuite) Test_CreateOrder_DbServiceWillIgnoreRequestStatus() {
	// ARRANGE
	json := `{
        "id": "input-entry-id",
        "name": "input-test-name",
        "medicineId": "input-test-medicine-id",
		"count": 20,
		"status": {
			"id": 5
		}
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("POST", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.CreateMedicineOrderEntry(ctx)

	// ASSERT
	var expectedObj = &MedicineOrderEntry{
		Id:         "input-entry-id",
		Name:       "input-test-name",
		MedicineId: "input-test-medicine-id",
		Count:      20,
		Status: Status{
			Id:               1,
			Value:            "To_ship",
			ValidTransitions: []int32{2, 4},
		},
	}
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			for _, entry := range arg.MedicineOrders {
				if reflect.DeepEqual(entry, *expectedObj) {
					return true
				}
			}
			return false
		}),
	)
}

func (suite *MedicineOrderSuite) Test_UpdateOrder_DbServiceUpdateSimpleFields() {
	// ARRANGE
	json := `{
		"count": 20
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.UpdateMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			for _, entry := range arg.MedicineOrders {
				if entry.Id == "test-entry" && entry.Count == 20 {
					return true
				}
			}
			return false
		}),
	)
}

func (suite *MedicineOrderSuite) Test_UpdateOrder_DbServiceCannotUpdateId() {
	// ARRANGE
	json := `{
		"id": "changed-entry-id"
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.UpdateMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertNotCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.Anything,
	)
}

func (suite *MedicineOrderSuite) Test_UpdateOrder_DbServiceUpdateStatus() {
	// ARRANGE
	json := `{
		"status": {
			"id": 2
		}
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.UpdateMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			for _, entry := range arg.MedicineOrders {
				if entry.Id == "test-entry" && entry.Status.Id == 2 {
					return true
				}
			}
			return false
		}),
	)
}

func (suite *MedicineOrderSuite) Test_UpdateOrder_DbServiceCannotUpdateInvalidStatus() {
	// ARRANGE
	json := `{
		"status": {
			"id": 3
		}
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.UpdateMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertNotCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.Anything,
	)
}

func (suite *MedicineOrderSuite) Test_UpdateOrder_DbServiceCannotChangeAnythingFromFinishedStatus() {
	// ARRANGE
	suite.dbAmbulanceServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Unset().
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id: "test-ambulance",
				MedicineOrders: []MedicineOrderEntry{
					{
						Id:         "test-entry",
						Name:       "test-name",
						MedicineId: "test-medicine-id",
						Count:      15,
						Status: Status{
							Id:               4,
							Value:            "Canceled",
							ValidTransitions: []int32{},
						},
					},
				},
			},
			nil,
		)

	json := `{
		"status": {
			"id": 2
		}
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_ambulance", suite.dbAmbulanceServiceMock)
	ctx.Set("db_service_status", suite.dbStatusServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/medicine-order/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineOrderAPI{}

	// ACT
	sut.UpdateMedicineOrderEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.dbAmbulanceServiceMock.AssertNotCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.Anything,
	)
}
