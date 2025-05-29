package medicine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/undy45/medicine-webapi/internal/db_service"
)

type MedicineSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Ambulance]
}

func TestMedicineSuite(t *testing.T) {
	suite.Run(t, new(MedicineSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *MedicineSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Ambulance]{}

	// Compile time Assert that the mock is of type db_service.DbService[Ambulance]
	var _ db_service.DbService[Ambulance] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id: "test-ambulance",
				MedicineInventory: []MedicineInventoryEntry{
					{
						Id:         "test-entry",
						Name:       "test-name",
						MedicineId: "test-medicine-id",
						Count:      15,
					},
				},
			},
			nil,
		)
}

func (suite *MedicineSuite) Test_DeleteInventory_DbService() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("DELETE", "/ambulance/test-ambulance/entries/test-entry", nil)

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.DeleteMedicineInventoryEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusNoContent, recorder.Code)
	suite.dbServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			return len(arg.MedicineInventory) == 0
		}),
	)
}

func (suite *MedicineSuite) Test_GetInventory_DbService() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("GET", "/ambulance/test-ambulance/entries/test-entry", nil)

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.GetMedicineInventoryEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	var respObj map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	suite.Equal("test-entry", respObj["id"])
	suite.Equal("test-name", respObj["name"])
	suite.Equal("test-medicine-id", respObj["medicineId"])
	suite.Equal(float64(15), respObj["count"])
	suite.dbServiceMock.AssertNotCalled(suite.T(), "UpdateDocument", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *MedicineSuite) Test_GetInventory_DbServiceGetAllEntries() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("GET", "/ambulance/test-ambulance/entries", nil)

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.GetMedicineInventoryEntries(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	var respObj []map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	var entry = respObj[0]
	suite.Equal("test-entry", entry["id"])
	suite.Equal("test-name", entry["name"])
	suite.Equal("test-medicine-id", entry["medicineId"])
	suite.Equal(float64(15), entry["count"])
	suite.dbServiceMock.AssertNotCalled(suite.T(), "UpdateDocument", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *MedicineSuite) Test_UpdateInventory_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
        "id": "test-entry",
        "name": "test-name",
        "medicineId": "test-medicine-id",
		"count": 20
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/ambulance/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.UpdateMedicineInventoryEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
}

func (suite *MedicineSuite) Test_UpdateInventory_DbServiceUpdateCalledWithCorrectInventory() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
        "id": "test-entry",
        "name": "test-name",
        "medicineId": "test-medicine-id",
		"count": 20
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/ambulance/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.UpdateMedicineInventoryEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			for _, entry := range arg.MedicineInventory {
				if entry.Id == "test-entry" && entry.Count == 20 {
					return true
				}
			}
			return false
		}),
	)
}

func (suite *MedicineSuite) Test_UpdateInventory_DbServiceDeleteInventoryWhenCountIsZero() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
        "id": "test-entry",
        "name": "test-name",
        "medicineId": "test-medicine-id",
		"count": 0
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/ambulance/test-ambulance/entries/test-entry", strings.NewReader(json))

	sut := implMedicineInventoryAPI{}

	// ACT
	sut.UpdateMedicineInventoryEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-ambulance",
		mock.MatchedBy(func(arg *Ambulance) bool {
			return len(arg.MedicineInventory) == 0
		}),
	)
}
