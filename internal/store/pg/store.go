package pg

import (
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Migrate driver.
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver.

	"github.com/imarrche/jwt-auth-example/internal/config"
	"github.com/imarrche/jwt-auth-example/internal/logger"
	"github.com/imarrche/jwt-auth-example/internal/store"
)

var (
	s    *Store
	once sync.Once
)

// Store is PostgreSQL store.
type Store struct {
	config   *config.PostgreSQL
	db       *sqlx.DB
	userRepo *userRepo
}

// Get creates store instance once and returns it.
func Get(c *config.PostgreSQL) *Store {
	once.Do(func() {
		s = &Store{config: c}
	})

	return s
}

// Open opens a connection with PostgreSQL.
func (s *Store) Open() error {
	// Connecting.
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.config.Host, s.config.Port, s.config.User,
		s.config.Password, s.config.DBName, s.config.SSLMode,
	)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	// Running migrations.
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/store/pg/migrations", "postgres", driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	logger.Get().Info("migrated PostgreSQL")

	s.db = db

	return nil
}

// Users returns the users repository.
func (s *Store) Users() store.UserRepo {
	if s.userRepo == nil {
		s.userRepo = newUserRepo(s.db)
	}

	return s.userRepo
}

// Close closes a connection with PostgreSQL.
func (s *Store) Close() error {
	return s.db.Close()
}
