package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
	foxRepository  *FoxRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		fmt.Println("Не открлась!")
		return err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Не пинганулась!")
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Fox() *FoxRepository {
	if s.foxRepository != nil {
		return s.foxRepository
	}

	s.foxRepository = &FoxRepository{
		store: s,
	}

	return s.foxRepository
}

type StoreError struct {
	textError string
}

func (se *StoreError) Error() string {
	return se.textError
}
