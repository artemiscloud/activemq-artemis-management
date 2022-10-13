package artemis

import (
	"fmt"
	"testing"

	"github.com/artemiscloud/activemq-artemis-management/jolokia"
	mock_jolokia "github.com/artemiscloud/activemq-artemis-management/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	j := mock_jolokia.NewMockIJolokia(ctrl)

	artemis := createMockArtemis(j)

	expectedStatus := "{\"properties\":{\"a_status.properties\": { \"cr:alder32\": \"3d8706a6\"}}}"
	j.
		EXPECT().
		Read(gomock.Eq("org.apache.activemq.artemis:broker=\\\"someBroker\\\"/Status")).
		DoAndReturn(func(_ string) (*jolokia.ResponseData, error) {
			return &jolokia.ResponseData{
				Status:    200,
				Value:     expectedStatus,
				ErrorType: "",
				Error:     "",
			}, nil
		}).
		AnyTimes()
	data, err := artemis.GetStatus()

	assert.Equal(t, expectedStatus, data)
	assert.Nil(t, err)
}

func TestGetStatusWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	j := mock_jolokia.NewMockIJolokia(ctrl)

	artemis := createMockArtemis(j)

	j.
		EXPECT().
		Read(gomock.Eq("org.apache.activemq.artemis:broker=\\\"someBroker\\\"/Status")).
		DoAndReturn(func(_ string) (*jolokia.ResponseData, error) {
			return &jolokia.ResponseData{
				Status:    404,
				Value:     "",
				ErrorType: "javax.management.AttributeNotFoundException",
				Error:     "javax.management.AttributeNotFoundException : No such attribute: Status",
			}, fmt.Errorf("javax.management.AttributeNotFoundException")
		}).
		AnyTimes()
	data, err := artemis.GetStatus()

	assert.Empty(t, data)
	assert.Error(t, err)
}

func createMockArtemis(j jolokia.IJolokia) Artemis {
	return Artemis{
		ip:          "0.0.0.0",
		jolokiaPort: "8161",
		name:        "someBroker",
		jolokia:     j,
	}
}
