package pg

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/imarrche/jwt-auth-example/internal/model"
)

func TestUserRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func([]model.User)
		expUsers []model.User
		expError bool
	}{
		{
			name: "users are retrieved",
			mock: func(us []model.User) {
				rows := sqlmock.NewRows([]string{"id", "username"})
				for _, u := range us {
					rows = rows.AddRow(u.ID, u.Username)
				}
				mock.ExpectQuery("SELECT (.+) FROM users;").WillReturnRows(rows)
			},
			expUsers: []model.User{
				{ID: 1, Username: "user1"}, {ID: 2, Username: "user2"},
			},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.expUsers)

		us, err := r.GetAll()

		if !tc.expError {
			assert.NoError(t, err)
			assert.Equal(t, tc.expUsers, us)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestUserRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func(model.User)
		user     model.User
		expUser  model.User
		expError bool
	}{
		{
			name: "user is created",
			mock: func(u model.User) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+) RETURNING id;").WithArgs(
					u.Username, u.Email, u.FirstName, u.SecondName, u.PasswordHash,
				).WillReturnRows(rows)
			},
			user:     model.User{Username: "user1"},
			expUser:  model.User{ID: 1, Username: "user1"},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.user)

		u, err := r.Create(tc.user)

		if !tc.expError {
			assert.NoError(t, err)
			assert.Equal(t, tc.expUser, u)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestUserRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func(model.User)
		user     model.User
		expUser  model.User
		expError bool
	}{
		{
			name: "user is retrieved by ID",
			mock: func(u model.User) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(
					u.ID, u.Username,
				)
				mock.ExpectQuery(
					"SELECT (.+) FROM users WHERE id = (.+);",
				).WithArgs(u.ID).WillReturnRows(rows)
			},
			user:     model.User{ID: 1, Username: "user1"},
			expUser:  model.User{ID: 1, Username: "user1"},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.user)

		u, err := r.GetByID(tc.user.ID)

		if !tc.expError {
			assert.NoError(t, err)
			assert.Equal(t, tc.expUser, u)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func(model.User)
		user     model.User
		expUser  model.User
		expError bool
	}{
		{
			name: "user is retrieved by email",
			mock: func(u model.User) {
				rows := sqlmock.NewRows([]string{"email", "username"}).AddRow(
					u.Email, u.Username,
				)
				mock.ExpectQuery(
					"SELECT (.+) FROM users WHERE email = (.+);",
				).WithArgs(u.Email).WillReturnRows(rows)
			},
			user:     model.User{Email: "user1@test.com", Username: "user1"},
			expUser:  model.User{Email: "user1@test.com", Username: "user1"},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.user)

		u, err := r.GetByEmail(tc.user.Email)

		if !tc.expError {
			assert.NoError(t, err)
			assert.Equal(t, tc.expUser, u)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestUserRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func(model.User)
		user     model.User
		expUser  model.User
		expError bool
	}{
		{
			name: "user is updated",
			mock: func(u model.User) {
				mock.ExpectExec(
					"UPDATE users SET (.+) WHERE id = (.+);",
				).WithArgs(
					u.Username, u.Email, u.FirstName, u.SecondName, u.PasswordHash, u.ID,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user:     model.User{ID: 1, Username: "updated_user1"},
			expUser:  model.User{ID: 1, Username: "updated_user1"},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.user)

		u, err := r.Update(tc.user)

		if !tc.expError {
			assert.NoError(t, err)
			assert.Equal(t, tc.expUser, u)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestUserRepo_DeleteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := newUserRepo(sqlx.NewDb(db, "postgres"))

	testcases := []struct {
		name     string
		mock     func(model.User)
		user     model.User
		expError bool
	}{
		{
			name: "user is deleted",
			mock: func(u model.User) {
				mock.ExpectExec(
					"DELETE FROM users WHERE id = (.+);",
				).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user:     model.User{ID: 1, Username: "user1"},
			expError: false,
		},
	}

	for _, tc := range testcases {
		tc.mock(tc.user)

		err := r.DeleteByID(tc.user.ID)

		if !tc.expError {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
