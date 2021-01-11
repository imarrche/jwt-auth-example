package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/imarrche/jwt-auth-example/internal/model"
	mock_service "github.com/imarrche/jwt-auth-example/internal/service/mocks"
	service_mock "github.com/imarrche/jwt-auth-example/internal/service/mocks"
)

func TestServer_signUp(t *testing.T) {
	server := &Server{router: chi.NewRouter()}
	server.configureRouter()

	testcases := []struct {
		name    string
		mock    func(*gomock.Controller, *mock_service.MockService, model.User)
		user    model.User
		expUser model.User
		expCode int
	}{
		{
			name: "user is signed up",
			mock: func(c *gomock.Controller, s *mock_service.MockService, u model.User) {
				as := mock_service.NewMockAuth(c)
				as.EXPECT().SignUp(u).Return(u, nil)
				s.EXPECT().Auth().Return(as)
			},
			user:    model.User{Username: "user1"},
			expUser: model.User{Username: "user1"},
			expCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		s := service_mock.NewMockService(c)
		tc.mock(c, s, tc.user)
		server.service = s

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.user)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", b)

		server.signUp().ServeHTTP(w, r)
		var u model.User
		err := json.NewDecoder(w.Body).Decode(&u)

		assert.NoError(t, err)
		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expUser, u)
	}
}

func TestServer_signIn(t *testing.T) {
	server := &Server{router: chi.NewRouter()}
	server.configureRouter()

	testcases := []struct {
		name        string
		mock        func(*gomock.Controller, *mock_service.MockService, signInRequest)
		request     signInRequest
		expResponse signInResponse
		expCode     int
	}{
		{
			name: "user is signed in",
			mock: func(c *gomock.Controller, s *mock_service.MockService, r signInRequest) {
				as := mock_service.NewMockAuth(c)
				as.EXPECT().SignIn(r.Email, r.Password).Return(
					"access_token", "refresh_token", nil,
				)
				s.EXPECT().Auth().Return(as)
			},
			request: signInRequest{
				Email: "user@test.com", Password: "password",
			},
			expResponse: signInResponse{
				AccessToken: "access_token", RefreshToken: "refresh_token",
			},
			expCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		s := service_mock.NewMockService(c)
		tc.mock(c, s, tc.request)
		server.service = s

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.request)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", b)

		server.signIn().ServeHTTP(w, r)
		var response signInResponse
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.NoError(t, err)
		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expResponse, response)
	}
}

func TestServer_refresh(t *testing.T) {
	server := &Server{router: chi.NewRouter()}
	server.configureRouter()

	testcases := []struct {
		name        string
		mock        func(*gomock.Controller, *mock_service.MockService, refreshRequest)
		request     refreshRequest
		expResponse refreshResponse
		expCode     int
	}{
		{
			name: "token is refreshed",
			mock: func(c *gomock.Controller, s *mock_service.MockService, r refreshRequest) {
				as := mock_service.NewMockAuth(c)
				as.EXPECT().RefreshAccessJWT(r.RefreshToken).Return("new_access_token", nil)
				s.EXPECT().Auth().Return(as)
			},
			request:     refreshRequest{RefreshToken: "refresh_token"},
			expResponse: refreshResponse{AccessToken: "new_access_token"},
			expCode:     http.StatusOK,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		s := service_mock.NewMockService(c)
		tc.mock(c, s, tc.request)
		server.service = s

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.request)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", b)

		server.refresh().ServeHTTP(w, r)
		var response refreshResponse
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.NoError(t, err)
		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expResponse, response)
	}
}

func TestServer_public(t *testing.T) {
	server := &Server{router: chi.NewRouter()}
	server.configureRouter()

	testcases := []struct {
		name    string
		expCode int
	}{
		{
			name:    "OK",
			expCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/public", nil)

		server.public().ServeHTTP(w, r)

		assert.Equal(t, tc.expCode, w.Code)
	}
}

func TestServer_private(t *testing.T) {
	server := &Server{router: chi.NewRouter()}
	server.configureRouter()

	testcases := []struct {
		name    string
		expCode int
	}{
		{
			name:    "OK",
			expCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/private", nil)

		server.private().ServeHTTP(w, r)

		assert.Equal(t, tc.expCode, w.Code)
	}
}
