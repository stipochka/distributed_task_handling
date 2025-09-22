package handler

import (
	"bytes"
	"fmt"
	"gateway/internal/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name         string
		taskType     string
		payload      string
		mockFunc     func(ms *mocks.TaskService)
		taskID       string
		respError    string
		expectedCode int
	}{
		{
			name:     "success",
			taskType: "resize image",

			mockFunc: func(ms *mocks.TaskService) {
				taskID := uuid.New()
				ms.
					On("SendTask", mock.Anything, mock.AnythingOfType("*models.Task")).
					Return(taskID, nil)
			},
			payload:      `{"image_url": "http://example.com/images/dsad.png", "height": 500, "width": 200}`,
			expectedCode: 201,
		},
		{
			name:     "service error",
			taskType: "resize image",
			payload:  `{"image_url": "http://example.com/images/dsad.png", "height": 500, "width": 200}`,
			mockFunc: func(ms *mocks.TaskService) {
				ms.
					On("SendTask", mock.Anything, mock.AnythingOfType("*models.Task")).
					Return(uuid.Nil, fmt.Errorf("some error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "empty task type",
			taskType:     "",
			payload:      `{"image_url": "http://example.com/images/dsad.png", "height": 500, "width": 200}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid request",
			taskType:     "ad121{}",
			payload:      `{"image_url": "http://example.com/images/dsad.pn "height": 500, "width": 200}`,
			expectedCode: 400,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			taskService := mocks.NewTaskService(t)

			if tc.mockFunc != nil {
				tc.mockFunc(taskService)
			}

			handler := NewHandler(taskService)

			input := fmt.Sprintf(`{"type":"%s", "payload": %s}`, tc.taskType, tc.payload)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/task", bytes.NewReader([]byte(input)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.InitRoute().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

		})
	}
}
