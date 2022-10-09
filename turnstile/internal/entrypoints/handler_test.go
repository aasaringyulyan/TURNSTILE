package entrypoints

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"turnstile/internal/models"
	"turnstile/internal/service"
	mock_service "turnstile/internal/service/mocks"
	"turnstile/pkg/logging"
)

const (
	debugLevel = "debug"
)

func TestHandler_getPassage(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTurnstile, check models.PassageCheck)

	tests := []struct {
		name                 string
		inputBody            string
		inputPassage         models.PassageCheck
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"key_hex": "5120AD2409C22FA2", "direction": 1}`,
			inputPassage: models.PassageCheck{
				KeyHex:    "5120AD2409C22FA2",
				Direction: 1,
			},
			mockBehavior: func(r *mock_service.MockTurnstile, check models.PassageCheck) {
				r.EXPECT().CheckHandler(check).Return(true, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accepted":true}`,
		},
		{
			name:                 "Wrong Input key_hex",
			inputBody:            `{"direction": 1}`,
			inputPassage:         models.PassageCheck{},
			mockBehavior:         func(r *mock_service.MockTurnstile, check models.PassageCheck) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:                 "Wrong Input direction",
			inputBody:            `{"key_hex": "5120AD2409C22FA2"}`,
			inputPassage:         models.PassageCheck{},
			mockBehavior:         func(r *mock_service.MockTurnstile, check models.PassageCheck) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"key_hex": "5120AD2409C22FA2", "direction": 1}`,
			inputPassage: models.PassageCheck{
				KeyHex:    "5120AD2409C22FA2",
				Direction: 1,
			},
			mockBehavior: func(r *mock_service.MockTurnstile, check models.PassageCheck) {
				r.EXPECT().CheckHandler(check).Return(false, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			logger := logging.GetLogger(debugLevel)

			turnstile := mock_service.NewMockTurnstile(c)
			test.mockBehavior(turnstile, test.inputPassage)

			tService := service.Turnstile(turnstile)
			handler := Handler{tService, logger}

			// Init Endpoint
			r := gin.New()
			r.POST("/check", handler.getPassage)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/check",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getLogs(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTurnstile, logs models.PassageLogsLinux)

	tests := []struct {
		name                 string
		inputBody            string
		inputLogs            models.PassageLogsLinux
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"logs": [{"logId": 0, "time": 1659568133, "direction": 1, "keyHex": "5120AD2409C22FA2"}]}`,
			inputLogs: models.PassageLogsLinux{
				Logs: []models.PassageLogLinux{
					{
						LogId:     0,
						Time:      1659568133,
						Direction: 1,
						Card:      "5120AD2409C22FA2",
					},
				},
			},
			mockBehavior: func(r *mock_service.MockTurnstile, logs models.PassageLogsLinux) {
				r.EXPECT().LogHandler(logs).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"code":200}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"logs": []}`,
			inputLogs:            models.PassageLogsLinux{},
			mockBehavior:         func(r *mock_service.MockTurnstile, logs models.PassageLogsLinux) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"logs": [{"logId": 0, "time": 1659568133, "direction": 1, "keyHex": "5120AD2409C22FA2"}]}`,
			inputLogs: models.PassageLogsLinux{
				Logs: []models.PassageLogLinux{
					{
						LogId:     0,
						Time:      1659568133,
						Direction: 1,
						Card:      "5120AD2409C22FA2",
					},
				},
			},
			mockBehavior: func(r *mock_service.MockTurnstile, logs models.PassageLogsLinux) {
				r.EXPECT().LogHandler(logs).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			logger := logging.GetLogger(debugLevel)

			turnstile := mock_service.NewMockTurnstile(c)
			test.mockBehavior(turnstile, test.inputLogs)

			tService := service.Turnstile(turnstile)
			handler := Handler{tService, logger}

			// Init Endpoint
			r := gin.New()
			r.POST("/logs", handler.getLogs)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/logs",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
