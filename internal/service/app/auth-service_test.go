package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/imarrche/jwt-auth-example/internal/model"
	mock_store "github.com/imarrche/jwt-auth-example/internal/store/mocks"
)

func TestAuthService_SignUp(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*gomock.Controller, *mock_store.MockStore, model.User)
		user     model.User
		expError bool
	}{
		{
			name: "user is signed up",
			mock: func(c *gomock.Controller, s *mock_store.MockStore, u model.User) {
				u.PasswordHash = "hashed_" + u.Password
				u.Password = ""

				ur := mock_store.NewMockUserRepo(c)
				ur.EXPECT().Create(gomock.Any()).Return(u, nil)
				s.EXPECT().Users().Return(ur)
			},
			user: model.User{
				ID:         1,
				Username:   "user1",
				Email:      "user1@test.com",
				FirstName:  "Name",
				SecondName: "Secondname",
				Password:   "password1",
			},
			expError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mock_store.NewMockStore(c)
			tc.mock(c, store, tc.user)
			s := newAuthService(store)
			u, err := s.SignUp(tc.user)

			if !tc.expError {
				assert.NoError(t, err)
				assert.Equal(t, tc.user.ID, u.ID)
				assert.Equal(t, "", u.Password)
				assert.NotEqual(t, "", u.PasswordHash)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAuthService_generateJWT(t *testing.T) {
	testcases := []struct {
		name     string
		userID   int
		expError bool
	}{
		{
			name:     "access JWT is generated",
			userID:   1,
			expError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			s := newAuthService(nil)
			token, err := s.generateJWT(tc.userID, "access")

			if !tc.expError {
				assert.NoError(t, err)
				assert.NotEqual(t, "", token)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAuthService_SignIn(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*gomock.Controller, *mock_store.MockStore, model.User)
		user     model.User
		expError bool
	}{
		{
			name: "user is signed in",
			mock: func(c *gomock.Controller, s *mock_store.MockStore, u model.User) {
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
				if err != nil {
					t.Fatal(err)
				}
				u.PasswordHash = string(hash)

				ur := mock_store.NewMockUserRepo(c)
				ur.EXPECT().GetByEmail(u.Email).Return(u, nil)
				s.EXPECT().Users().Return(ur)
			},
			user: model.User{
				ID:       1,
				Email:    "user1@test.com",
				Password: "password1",
			},
			expError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mock_store.NewMockStore(c)
			tc.mock(c, store, tc.user)
			s := newAuthService(store)
			accessJWT, refreshJWT, err := s.SignIn(tc.user.Email, tc.user.Password)

			if !tc.expError {
				assert.NoError(t, err)
				assert.NotEqual(t, "", accessJWT)
				assert.NotEqual(t, "", refreshJWT)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAuthService_ValidateJWT(t *testing.T) {
	testcases := []struct {
		name      string
		userID    int
		tokenType string
		expError  bool
	}{
		{
			name:     "token is valid",
			userID:   1,
			expError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			s := newAuthService(nil)
			token, err := s.generateJWT(tc.userID, tc.tokenType)
			if err != nil {
				t.Fatal(err)
			}
			userID, err := s.ValidateJWT(token, tc.tokenType)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.expError {
				assert.NoError(t, err)
				assert.Equal(t, tc.userID, userID)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAuthService_RefreshAccessJWT(t *testing.T) {
	testcases := []struct {
		name     string
		userID   int
		expError bool
	}{
		{
			name:     "access token is refreshed",
			userID:   1,
			expError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			s := newAuthService(nil)
			token, err := s.generateJWT(tc.userID, "refresh")
			if err != nil {
				t.Fatal(err)
			}
			accessJWT, err := s.RefreshAccessJWT(token)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.expError {
				assert.NoError(t, err)
				assert.NotEqual(t, "", accessJWT)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
