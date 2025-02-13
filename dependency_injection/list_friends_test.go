package friends

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chr1sbest/api.mobl.ai/internal/auth"
	"github.com/chr1sbest/api.mobl.ai/internal/storage/model"
	"github.com/chr1sbest/api.mobl.ai/internal/util"
	"github.com/chr1sbest/api.mobl.ai/internal/util/errors"

	authMock "github.com/chr1sbest/api.mobl.ai/internal/auth/mocks"
	storageMock "github.com/chr1sbest/api.mobl.ai/internal/storage/mocks"
)

func setupListFriendsTests() (*storageMock.Storage, *authMock.Authenticator, *ListFriendsHandler) {
	mockStorage := new(storageMock.Storage)
	mockAuth := new(authMock.Authenticator)
	handler := &ListFriendsHandler{
		Db:   mockStorage,
		Auth: mockAuth,
	}
	return mockStorage, mockAuth, handler
}

func TestHandleListFriends(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		request       ListFriendsRequest
		setupMocks    func(ms *storageMock.Storage, ma *authMock.Authenticator)
		expectedCode  int
		expectedError error
		expectedBody  []*model.Friend
	}{
		{
			name: "Successful List Friends",
			request: ListFriendsRequest{
				Token: "validToken",
			},
			setupMocks: func(ms *storageMock.Storage, ma *authMock.Authenticator) {
				ma.On("ValidateIDToken", ctx, mock.Anything).Return(&auth.CustomClaims{Email: "user@example.com"}, nil)
				ms.On("ListFriends", ctx, "user@example.com").Return([]*model.Friend{
					{
						ID:               "friend1",
						Email:            "user@example.com",
						MirrorID:         "mirror1",
						FriendName:       "John Doe",
						FriendStreak:     5,
						FriendEXP:        100,
						FriendLastActive: "2024-09-01T12:00:00Z",
						FriendCreatedAt:  "2024-08-01T12:00:00Z",
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: []*model.Friend{
				{
					ID:               "friend1",
					Email:            "user@example.com",
					MirrorID:         "mirror1",
					FriendName:       "John Doe",
					FriendStreak:     5,
					FriendEXP:        100,
					FriendLastActive: "2024-09-01T12:00:00Z",
					FriendCreatedAt:  "2024-08-01T12:00:00Z",
				},
			},
		},
		{
			name: "Missing Token in Request",
			request: ListFriendsRequest{
				Token: "",
			},
			setupMocks:    func(ms *storageMock.Storage, ma *authMock.Authenticator) {},
			expectedCode:  http.StatusBadRequest,
			expectedError: errors.MissingUserToken,
		},
		{
			name: "ID Token Validation Fails",
			request: ListFriendsRequest{
				Token: "invalidToken",
			},
			setupMocks: func(ms *storageMock.Storage, ma *authMock.Authenticator) {
				ma.On("ValidateIDToken", ctx, mock.Anything).Return(nil, errors.IDTokenInvalid)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: errors.IDTokenInvalid,
		},
		{
			name: "Database Operation Fails",
			request: ListFriendsRequest{
				Token: "validToken",
			},
			setupMocks: func(ms *storageMock.Storage, ma *authMock.Authenticator) {
				ma.On("ValidateIDToken", ctx, mock.Anything).Return(&auth.CustomClaims{Email: "user@example.com"}, nil)
				ms.On("ListFriends", ctx, "user@example.com").Return(nil, errors.FailedFriendList)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.FailedFriendList,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStorage, mockAuth, handler := setupListFriendsTests()
			tc.setupMocks(mockStorage, mockAuth)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/friends/list", serializeRequest(tc.request))
			resp := httptest.NewRecorder()

			handler.Handle(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)

			if tc.expectedError != nil {
				responseBody := resp.Body.String()
				util.AssertErrorResponse(t, responseBody, tc.expectedError)
			} else {
				var actualResponse []*model.Friend
				err := json.NewDecoder(resp.Body).Decode(&actualResponse)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualResponse)
			}

			mockStorage.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}
