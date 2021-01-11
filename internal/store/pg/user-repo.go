package pg

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/imarrche/jwt-auth-example/internal/model"
	"github.com/imarrche/jwt-auth-example/internal/store"
)

// userRepo is the user repository for PostgreSQL store.
type userRepo struct {
	db *sqlx.DB
}

// newUserRepo creates and returns a new userRepo instance.
func newUserRepo(db *sqlx.DB) *userRepo { return &userRepo{db: db} }

// GetAll returns all users.
func (r *userRepo) GetAll() ([]model.User, error) {
	users := []model.User{}
	if err := r.db.Select(&users, "SELECT * FROM users;"); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

// Create creates and returns a new user.
func (r *userRepo) Create(u model.User) (model.User, error) {
	query := "INSERT INTO users (username, email, first_name, second_name, password_hash) "
	query += "VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	row := r.db.QueryRow(query, u.Username, u.Email, u.FirstName, u.SecondName, u.PasswordHash)

	var id int
	err := row.Scan(&id)
	if err, ok := err.(*pq.Error); ok {
		if err.Constraint == "users_username_key" {
			return model.User{}, store.ErrUsernameIsTaken
		} else if err.Constraint == "users_email_key" {
			return model.User{}, store.ErrEmailIsTaken
		}
	} else if err != nil {
		return model.User{}, err
	}
	u.ID = id

	return u, nil
}

// GetByID returns the user with specific ID.
func (r *userRepo) GetByID(id int) (model.User, error) {
	u := model.User{}
	if err := r.db.Get(&u, "SELECT * FROM users WHERE id = $1;", id); err != nil {
		return model.User{}, err
	}

	return u, nil
}

// GetByEmail returns the user with specific email.
func (r *userRepo) GetByEmail(email string) (model.User, error) {
	u := model.User{}
	if err := r.db.Get(&u, "SELECT * FROM users WHERE email = $1;", email); err != nil {
		return model.User{}, err
	}

	return u, nil
}

// Update updates the user.
func (r *userRepo) Update(u model.User) (model.User, error) {
	query := "UPDATE users SET username = :username, email = :email, first_name = :first_name, "
	query += "second_name = :second_name, password_hash = :password_hash WHERE id = :id;"
	res, err := r.db.NamedExec(query, u)
	if err != nil {
		return model.User{}, err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return model.User{}, err
	} else if rowsCount == 0 {
		return model.User{}, errors.New("not found")
	}

	return u, nil
}

// DeleteByID deletes the user with specific ID.
func (r *userRepo) DeleteByID(id int) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsCount == 0 {
		return errors.New("not found")
	}

	return nil
}
