package models

import (
	"database/sql"
	"time"
)

type Users struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return -1, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
