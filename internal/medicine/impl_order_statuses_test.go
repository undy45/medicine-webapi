package medicine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/undy45/medicine-webapi/internal/db_service"
)

type OrderStatusesSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Status]
}

func TestOrderStatusesSuite(t *testing.T) {
	suite.Run(t, new(OrderStatusesSuite))
}

func (suite *OrderStatusesSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Status]{}

	// Compile time Assert that the mock is of type db_service.DbService[Status]
	var _ db_service.DbService[Status] = suite.dbServiceMock

	suite.dbServiceMock.
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
	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Status{
				Id:               1,
				Value:            "To_ship",
				ValidTransitions: []int32{2, 4},
			},
			nil,
		)
}

func (suite *OrderStatusesSuite) Test_GetStatus_DbService() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_status", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "statusId", Value: "1"},
	}
	ctx.Request = httptest.NewRequest("GET", "/medicine-order/statuses/1", nil)

	sut := implOrderStatusesApi{}

	// ACT
	sut.GetStatus(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	var respObj map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	checkStatus(suite, respObj, &Status{
		Id:               1,
		Value:            "To_ship",
		ValidTransitions: []int32{2, 4},
	})
}

func (suite *OrderStatusesSuite) Test_GetStatuses_DbService() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_status", suite.dbServiceMock)
	ctx.Request = httptest.NewRequest("GET", "/medicine-order/statuses", nil)

	sut := implOrderStatusesApi{}

	// ACT
	sut.GetStatuses(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	var respObj []map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &respObj)
	suite.Require().NoError(err)
	for _, gottenStatus := range respObj {
		switch gottenStatus["value"] {
		case "To_ship":
			checkStatus(suite, gottenStatus, &Status{
				Id:               1,
				Value:            "To_ship",
				ValidTransitions: []int32{2, 4},
			})
			continue
		case "Shipped":
			checkStatus(suite, gottenStatus, &Status{
				Id:               2,
				Value:            "Shipped",
				ValidTransitions: []int32{3, 4},
			})
			continue
		case "Delivered":
			checkStatus(suite, gottenStatus, &Status{
				Id:               3,
				Value:            "Delivered",
				ValidTransitions: []int32{},
			})
			continue
		case "Canceled":
			checkStatus(suite, gottenStatus, &Status{
				Id:               4,
				Value:            "Canceled",
				ValidTransitions: []int32{},
			})
			continue
		default:
			suite.Fail("Unexpected status value: " + gottenStatus["value"].(string))
		}
	}

}

func checkStatus(suite *OrderStatusesSuite, gottenStatus map[string]interface{}, expectedStatus *Status) {
	suite.Equal(expectedStatus.Id, int32(gottenStatus["id"].(float64)))
	suite.Equal(expectedStatus.Value, gottenStatus["value"])
	if len(expectedStatus.ValidTransitions) == 0 {
		suite.Nil(gottenStatus["valid_transitions"])
		return
	}
	validTransitions, ok := gottenStatus["validTransitions"].([]interface{})
	suite.True(ok)
	transitions := make([]int32, len(validTransitions))
	for i, v := range validTransitions {
		transitions[i] = int32(v.(float64))
	}
	suite.ElementsMatch(expectedStatus.ValidTransitions, transitions)
}
