package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/netbirdio/netbird/management/server"
	"github.com/netbirdio/netbird/management/server/http/api"
	"github.com/netbirdio/netbird/management/server/jwtclaims"
	"github.com/netbirdio/netbird/management/server/mock_server"
	"github.com/netbirdio/netbird/management/server/posture"
	"github.com/netbirdio/netbird/management/server/status"
)

func initPostureChecksTestData(postureChecks ...*posture.Checks) *PostureChecksHandler {
	testPostureChecks := make(map[string]*posture.Checks, len(postureChecks))
	for _, postureCheck := range postureChecks {
		testPostureChecks[postureCheck.ID] = postureCheck
	}

	return &PostureChecksHandler{
		accountManager: &mock_server.MockAccountManager{
			GetPostureChecksFunc: func(accountID, postureChecksID, userID string) (*posture.Checks, error) {
				p, ok := testPostureChecks[postureChecksID]
				if !ok {
					return nil, status.Errorf(status.NotFound, "posture checks not found")
				}
				return p, nil
			},
			SavePostureChecksFunc: func(accountID, userID string, postureChecks *posture.Checks) error {
				postureChecks.ID = "postureCheck"
				testPostureChecks[postureChecks.ID] = postureChecks
				return nil
			},
			DeletePostureChecksFunc: func(accountID, postureChecksID, userID string) error {
				_, ok := testPostureChecks[postureChecksID]
				if !ok {
					return status.Errorf(status.NotFound, "posture checks not found")
				}
				delete(testPostureChecks, postureChecksID)

				return nil
			},
			ListPostureChecksFunc: func(accountID, userID string) ([]*posture.Checks, error) {
				accountPostureChecks := make([]*posture.Checks, len(testPostureChecks))
				for _, p := range testPostureChecks {
					accountPostureChecks = append(accountPostureChecks, p)
				}
				return accountPostureChecks, nil
			},
			GetAccountFromTokenFunc: func(claims jwtclaims.AuthorizationClaims) (*server.Account, *server.User, error) {
				user := server.NewAdminUser("test_user")
				return &server.Account{
					Id:     claims.AccountId,
					Domain: "hotmail.com",
					Policies: []*server.Policy{
						{ID: "id-existed"},
					},
					Groups: map[string]*server.Group{
						"F": {ID: "F"},
						"G": {ID: "G"},
					},
					Users: map[string]*server.User{
						"test_user": user,
					},
				}, user, nil
			},
		},
		claimsExtractor: jwtclaims.NewClaimsExtractor(
			jwtclaims.WithFromRequestContext(func(r *http.Request) jwtclaims.AuthorizationClaims {
				return jwtclaims.AuthorizationClaims{
					UserId:    "test_user",
					Domain:    "hotmail.com",
					AccountId: "test_id",
				}
			}),
		),
	}
}

func TestGetPostureCheck(t *testing.T) {
	tt := []struct {
		name           string
		expectedStatus int
		expectedBody   bool
		requestType    string
		requestPath    string
		requestBody    io.Reader
	}{
		{
			name:           "GetPostureCheck OK",
			expectedBody:   true,
			requestType:    http.MethodGet,
			requestPath:    "/api/posture-checks/postureCheck",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetPostureCheck Not Found",
			requestType:    http.MethodGet,
			requestPath:    "/api/posture-checks/not-exists",
			expectedStatus: http.StatusNotFound,
		},
	}

	postureCheck := &posture.Checks{
		ID:   "postureCheck",
		Name: "name",
		Checks: []posture.Check{
			&posture.NBVersionCheck{
				Enabled:    true,
				MinVersion: "1.0.0",
			},
		},
	}

	p := initPostureChecksTestData(postureCheck)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tc.requestType, tc.requestPath, tc.requestBody)

			router := mux.NewRouter()
			router.HandleFunc("/api/posture-checks/{postureCheckId}", p.GetPostureCheck).Methods("GET")
			router.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
				return
			}

			if !tc.expectedBody {
				return
			}

			content, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("I don't know what I expected; %v", err)
			}

			var got api.PostureCheck
			if err = json.Unmarshal(content, &got); err != nil {
				t.Fatalf("Sent content is not in correct json format; %v", err)
			}

			assert.Equal(t, got.Id, postureCheck.ID)
			assert.Equal(t, got.Name, postureCheck.Name)
		})
	}
}

func TestPostureCheckUpdate(t *testing.T) {
	str := func(s string) *string { return &s }
	tt := []struct {
		name                 string
		expectedStatus       int
		expectedBody         bool
		expectedPostureCheck *api.PostureCheck
		requestType          string
		requestPath          string
		requestBody          io.Reader
	}{
		{
			name:        "Create Posture Checks",
			requestType: http.MethodPost,
			requestPath: "/api/posture-checks",
			requestBody: bytes.NewBuffer(
				[]byte(`{
		           "name": "default",
                  "description": "default",
		           "checks": {
						"nb_version_check": {
							"enabled": true,
							"min_version": "1.2.3",
							"max_version": "2.0.0"
		           		}
                  }
				}`)),
			expectedStatus: http.StatusOK,
			expectedBody:   true,
			expectedPostureCheck: &api.PostureCheck{
				Id:          "postureCheck",
				Name:        "default",
				Description: str("default"),
				Checks: &api.Checks{
					NbVersionCheck: &api.NBVersionCheck{
						Enabled:    true,
						MinVersion: "1.2.3",
						MaxVersion: str("2.0.0"),
					},
				},
			},
		},
		{
			name:        "Create Posture Checks Invalid Name",
			requestType: http.MethodPost,
			requestPath: "/api/posture-checks",
			requestBody: bytes.NewBuffer(
				[]byte(`{
                   "checks": {
						"nb_version_check": {
							"enabled": true,
							"min_version": "1.2.0"
                   	}
					}
				}`)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:        "Create Posture Checks Invalid NetBird's Min Version",
			requestType: http.MethodPost,
			requestPath: "/api/posture-checks",
			requestBody: bytes.NewBuffer(
				[]byte(`{
					"name": "default",
                   "checks": {
						"nb_version_check": {
							"enabled": true,
                   	}
					}
				}`)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:        "Update Posture Checks",
			requestType: http.MethodPut,
			requestPath: "/api/posture-checks/postureCheck",
			requestBody: bytes.NewBuffer(
				[]byte(`{
		           "name": "default",
		           "checks": {
						"nb_version_check": {
							"enabled": true,
							"min_version": "1.9.0"
		           		}
					}
				}`)),
			expectedStatus: http.StatusOK,
			expectedBody:   true,
			expectedPostureCheck: &api.PostureCheck{
				Id:          "postureCheck",
				Name:        "default",
				Description: str(""),
				Checks: &api.Checks{
					NbVersionCheck: &api.NBVersionCheck{
						Enabled:    true,
						MinVersion: "1.9.0",
						MaxVersion: str(""),
					},
				},
			},
		},
		{
			name:        "Update Posture Checks Invalid Name",
			requestType: http.MethodPut,
			requestPath: "/api/posture-checks/postureCheck",
			requestBody: bytes.NewBuffer(
				[]byte(`{
                   "checks": {
						"nb_version_check": {
							"enabled": true,
                   	}
					}
				}`)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name:        "Update Posture Checks Invalid NetBird's Min Version",
			requestType: http.MethodPut,
			requestPath: "/api/posture-checks/postureCheck",
			requestBody: bytes.NewBuffer(
				[]byte(`{
					"name": "default",
                   "checks": {
						"nb_version_check": {
							"enabled": false,
                   	}
					}
				}`)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
	}

	p := initPostureChecksTestData(&posture.Checks{
		ID:   "postureCheck",
		Name: "postureCheck",
		Checks: []posture.Check{
			&posture.NBVersionCheck{
				Enabled:    true,
				MinVersion: "1.0.0",
			},
		},
	})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tc.requestType, tc.requestPath, tc.requestBody)

			router := mux.NewRouter()
			router.HandleFunc("/api/posture-checks", p.CreatePostureCheck).Methods("POST")
			router.HandleFunc("/api/posture-checks/{postureCheckId}", p.UpdatePostureCheck).Methods("PUT")
			router.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			content, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("I don't know what I expected; %v", err)
				return
			}

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v, content: %s",
					status, tc.expectedStatus, string(content))
				return
			}

			if !tc.expectedBody {
				return
			}

			expected, err := json.Marshal(tc.expectedPostureCheck)
			if err != nil {
				t.Fatalf("marshal expected posture check: %v", err)
				return
			}

			assert.Equal(t, strings.Trim(string(content), " \n"), string(expected), "content mismatch")
		})
	}
}